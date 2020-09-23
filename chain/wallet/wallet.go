package wallet

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	cbg "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"

	"github.com/filecoin-project/go-state-types/crypto"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-address"

	_ "github.com/filecoin-project/lotus/lib/sigs/bls"  // enable bls signatures
	_ "github.com/filecoin-project/lotus/lib/sigs/secp" // enable secp signatures

	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/lib/sigs"
)

var log = logging.Logger("wallet")

const (
	KNamePrefix  = "wallet-"
	KTrashPrefix = "trash-"
	KDefault     = "default"
	KTBLS        = "bls"
	KTSecp256k1  = "secp256k1"
)

type Wallet struct {
	keys     map[address.Address]*Key
	keystore types.KeyStore

	lk sync.Mutex
}

func NewWallet(keystore types.KeyStore) (*Wallet, error) {
	w := &Wallet{
		keys:     make(map[address.Address]*Key),
		keystore: keystore,
	}

	return w, nil
}

func KeyWallet(keys ...*Key) *Wallet {
	m := make(map[address.Address]*Key)
	for _, key := range keys {
		m[key.Address] = key
	}

	return &Wallet{
		keys: m,
	}
}

func (w *Wallet) Sign(ctx context.Context, addr address.Address, msg []byte) (*crypto.Signature, error) {
	ki, err := w.findKey(addr)
	if err != nil {
		return nil, err
	}
	if ki == nil {
		return nil, xerrors.Errorf("signing using key '%s': %w", addr.String(), types.ErrKeyInfoNotFound)
	}

	return sigs.Sign(ActSigType(ki.Type), ki.PrivateKey, msg)
}

func (w *Wallet) findKey(addr address.Address) (*Key, error) {
	w.lk.Lock()
	defer w.lk.Unlock()

	k, ok := w.keys[addr]
	if ok {
		return k, nil
	}
	if w.keystore == nil {
		log.Warn("findKey didn't find the key in in-memory wallet")
		return nil, nil
	}

	ki, err := w.keystore.Get(KNamePrefix + addr.String())
	if err != nil {
		if xerrors.Is(err, types.ErrKeyInfoNotFound) {
			return nil, nil
		}
		return nil, xerrors.Errorf("getting from keystore: %w", err)
	}
	k, err = NewKey(ki)
	if err != nil {
		return nil, xerrors.Errorf("decoding from keystore: %w", err)
	}
	w.keys[k.Address] = k
	return k, nil
}

func (w *Wallet) Export(addr address.Address) (*types.KeyInfo, error) {
	k, err := w.findKey(addr)
	if err != nil {
		return nil, xerrors.Errorf("failed to find key to export: %w", err)
	}

	return &k.KeyInfo, nil
}

func (w *Wallet) Import(ki *types.KeyInfo) (address.Address, error) {
	w.lk.Lock()
	defer w.lk.Unlock()

	k, err := NewKey(*ki)
	if err != nil {
		return address.Undef, xerrors.Errorf("failed to make key: %w", err)
	}

	if err := w.keystore.Put(KNamePrefix+k.Address.String(), k.KeyInfo); err != nil {
		return address.Undef, xerrors.Errorf("saving to keystore: %w", err)
	}

	return k.Address, nil
}

func (w *Wallet) ListAddrs() ([]address.Address, error) {
	all, err := w.keystore.List()
	if err != nil {
		return nil, xerrors.Errorf("listing keystore: %w", err)
	}

	sort.Strings(all)

	out := make([]address.Address, 0, len(all))
	for _, a := range all {
		if strings.HasPrefix(a, KNamePrefix) {
			name := strings.TrimPrefix(a, KNamePrefix)
			addr, err := address.NewFromString(name)
			if err != nil {
				return nil, xerrors.Errorf("converting name to address: %w", err)
			}
			out = append(out, addr)
		}
	}

	return out, nil
}

func (w *Wallet) GetDefault() (address.Address, error) {
	w.lk.Lock()
	defer w.lk.Unlock()

	ki, err := w.keystore.Get(KDefault)
	if err != nil {
		return address.Undef, xerrors.Errorf("failed to get default key: %w", err)
	}

	k, err := NewKey(ki)
	if err != nil {
		return address.Undef, xerrors.Errorf("failed to read default key from keystore: %w", err)
	}

	return k.Address, nil
}

func (w *Wallet) SetDefault(a address.Address) error {
	w.lk.Lock()
	defer w.lk.Unlock()

	ki, err := w.keystore.Get(KNamePrefix + a.String())
	if err != nil {
		return err
	}

	if err := w.keystore.Delete(KDefault); err != nil {
		if !xerrors.Is(err, types.ErrKeyInfoNotFound) {
			log.Warnf("failed to unregister current default key: %s", err)
		}
	}

	if err := w.keystore.Put(KDefault, ki); err != nil {
		return err
	}

	return nil
}

func GenerateKey(typ crypto.SigType) (*Key, error) {
	pk, err := sigs.Generate(typ)
	if err != nil {
		return nil, err
	}
	ki := types.KeyInfo{
		Type:       kstoreSigType(typ),
		PrivateKey: pk,
	}
	return NewKey(ki)
}

func (w *Wallet) GenerateKey(typ crypto.SigType) (address.Address, error) {
	w.lk.Lock()
	defer w.lk.Unlock()

	k, err := GenerateKey(typ)
	if err != nil {
		return address.Undef, err
	}

	if err := w.keystore.Put(KNamePrefix+k.Address.String(), k.KeyInfo); err != nil {
		return address.Undef, xerrors.Errorf("saving to keystore: %w", err)
	}
	w.keys[k.Address] = k

	_, err = w.keystore.Get(KDefault)
	if err != nil {
		if !xerrors.Is(err, types.ErrKeyInfoNotFound) {
			return address.Undef, err
		}

		if err := w.keystore.Put(KDefault, k.KeyInfo); err != nil {
			return address.Undef, xerrors.Errorf("failed to set new key as default: %w", err)
		}
	}

	return k.Address, nil
}

func (w *Wallet) HasKey(addr address.Address) (bool, error) {
	k, err := w.findKey(addr)
	if err != nil {
		return false, err
	}
	return k != nil, nil
}

func (w *Wallet) DeleteKey(addr address.Address) error {
	k, err := w.findKey(addr)
	if err != nil {
		return xerrors.Errorf("failed to delete key %s : %w", addr, err)
	}

	if err := w.keystore.Put(KTrashPrefix+k.Address.String(), k.KeyInfo); err != nil {
		return xerrors.Errorf("failed to mark key %s as trashed: %w", addr, err)
	}

	if err := w.keystore.Delete(KNamePrefix + k.Address.String()); err != nil {
		return xerrors.Errorf("failed to delete key %s: %w", addr, err)
	}

	return nil
}

type Key struct {
	types.KeyInfo

	PublicKey []byte
	Address   address.Address
}

func NewKey(keyinfo types.KeyInfo) (*Key, error) {
	k := &Key{
		KeyInfo: keyinfo,
	}

	var err error
	k.PublicKey, err = sigs.ToPublic(ActSigType(k.Type), k.PrivateKey)
	if err != nil {
		return nil, err
	}

	switch k.Type {
	case KTSecp256k1:
		k.Address, err = address.NewSecp256k1Address(k.PublicKey)
		if err != nil {
			return nil, xerrors.Errorf("converting Secp256k1 to address: %w", err)
		}
	case KTBLS:
		k.Address, err = address.NewBLSAddress(k.PublicKey)
		if err != nil {
			return nil, xerrors.Errorf("converting BLS to address: %w", err)
		}
	default:
		return nil, xerrors.Errorf("unknown key type")
	}
	return k, nil

}

func kstoreSigType(typ crypto.SigType) string {
	switch typ {
	case crypto.SigTypeBLS:
		return KTBLS
	case crypto.SigTypeSecp256k1:
		return KTSecp256k1
	default:
		return ""
	}
}

func ActSigType(typ string) crypto.SigType {
	switch typ {
	case KTBLS:
		return crypto.SigTypeBLS
	case KTSecp256k1:
		return crypto.SigTypeSecp256k1
	default:
		return 0
	}
}

const (
	dsKeyAddrNonce = "AddressNonce"
)

type MessageSigner struct {
	*Wallet
	ds datastore.Batching
}

func NewMessageSigner(wallet *Wallet, ds dtypes.MetadataDS) *MessageSigner {
	ds = namespace.Wrap(ds, datastore.NewKey("/message-signer/"))
	return &MessageSigner{
		Wallet: wallet,
		ds:     ds,
	}
}

func (ms *MessageSigner) SignMessage(ctx context.Context, msg *types.Message) (*types.SignedMessage, error) {
	nonce, err := ms.nextNonce(msg.From)
	if err != nil {
		return nil, xerrors.Errorf("failed to create nonce: %w", err)
	}

	msg.Nonce = nonce
	sig, err := ms.Sign(ctx, msg.From, msg.Cid().Bytes())
	if err != nil {
		return nil, xerrors.Errorf("failed to sign message: %w", err)
	}

	return &types.SignedMessage{
		Message:   *msg,
		Signature: *sig,
	}, nil
}

func (ms *MessageSigner) nextNonce(addr address.Address) (uint64, error) {
	addrNonceKey := datastore.KeyWithNamespaces([]string{dsKeyAddrNonce, addr.String()})

	// Get the nonce for this address from the datastore
	nonce := uint64(0)
	nonceBytes, err := ms.ds.Get(addrNonceKey)
	switch {
	case xerrors.Is(err, datastore.ErrNotFound):
		// No nonce yet for this address so just use zero
	case err != nil:
		return 0, xerrors.Errorf("failed to get nonce from datastore: %w", err)
	default:
		// There is a nonce already, so get it and increment
		maj, val, err := cbg.CborReadHeader(bytes.NewReader(nonceBytes))
		if err != nil {
			return 0, err
		}

		if maj != cbg.MajUnsignedInt {
			return 0, fmt.Errorf("bad cbor type")
		}

		nonce = val + 1
	}

	// Write the nonce to the datastore for this address
	buf := bytes.Buffer{}
	_, err = buf.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, nonce))
	if err != nil {
		return 0, xerrors.Errorf("failed to marshall nonce: %w", err)
	}
	err = ms.ds.Put(addrNonceKey, buf.Bytes())
	if err != nil {
		return 0, xerrors.Errorf("failed to write nonce to datastore: %w", err)
	}

	return nonce, nil
}
