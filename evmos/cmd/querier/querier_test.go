package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var querier *Querier
var _ = BeforeSuite(func() {
	var err error
	querier, err = NewQuerier(RPC_URL)
	Expect(err).To(BeNil())

	querier.storeBlocks(100, 102)
})

var _ = Describe("Querier", func() {
	Describe("storeBlocks", func() {
		Context("When storing blocks from 100 to 102", func() {
			It("should store 101 blocks in querier.blocks", func() {
				Expect(querier.blocks).To(HaveLen(3))
			})
		})
	})
	Describe("storeWallet", func() {
		It("should add non-contract addresses to walletBalances", func() {
			walletCache := make(map[string]struct{})
			initialWalletCount := len(querier.walletBalances)

			for _, block := range querier.blocks {
				for _, tx := range block.Transactions() {
					err := querier.storeWallet(tx, walletCache)
					Expect(err).To(BeNil())
				}
			}

			Expect(len(querier.walletBalances)).To(BeNumerically(">", initialWalletCount))
			Expect(len(walletCache)).To(BeNumerically(">", 0))
		})
	})

	Describe("populateWalletBalances", func() {
		It("should have at least one wallet with more than 0 balance", func() {
			walletCache := make(map[string]struct{})
			for _, block := range querier.blocks {
				for _, tx := range block.Transactions() {
					querier.storeWallet(tx, walletCache)
				}
			}

			err := querier.populateWalletBalances()
			Expect(err).To(BeNil())

			populated := false
			for _, wb := range querier.walletBalances {
				if wb.Balance.Cmp(common.Big0) > 0 {
					populated = true
					break
				}
			}

			Expect(populated).To(BeTrue(), "There should exist a wallet with positive balance.")
		})
	})

	Describe("Sorting Wallet Balances", func() {
		Context("When sorting wallets by balance", func() {
			It("should sort wallet balances from highest to lowest", func() {
				querier.walletBalances = []WalletBalance{
					{Address: "walletA", Balance: big.NewInt(10)},
					{Address: "walletB", Balance: big.NewInt(28)},
					{Address: "walletC", Balance: big.NewInt(5)},
				}

				querier.sortWalletsByBalance()

				Expect(querier.walletBalances[0].Balance).To(Equal(big.NewInt(28)))
				Expect(querier.walletBalances[1].Balance).To(Equal(big.NewInt(10)))
				Expect(querier.walletBalances[2].Balance).To(Equal(big.NewInt(5)))
			})
		})
	})

	Describe("storeContracts", func() {
		Context("When storing contract interactions", func() {
			It("should correctly store contract addresses from the transaction trace", func() {
				txHash := common.HexToHash("0x28bfa2817bffc06dab8aa8bc2b5524de32bd4296932b25132e86593bc368a8d4")

				err := querier.storeContracts(txHash)
				Expect(err).To(BeNil())

				Expect(len(querier.contracts)).To(BeNumerically(">", 0))
			})
		})
	})

	Describe("sortContractsByInteractions", func() {
		Context("When sorting contracts by interactions", func() {
			It("should return contracts sorted by interaction count in descending order", func() {
				querier.contracts = map[string]int{
					"0x1": 5,
					"0x2": 10,
					"0x3": 3,
				}

				sortedContracts := querier.sortContractsByInteractions()

				Expect(sortedContracts).To(HaveLen(3))
				Expect(sortedContracts[0]).To(Equal("0x2"))
				Expect(sortedContracts[1]).To(Equal("0x1"))
				Expect(sortedContracts[2]).To(Equal("0x3"))
			})
		})
	})

	Describe("getTransactionTrace", func() {
		Context("When getting a transaction trace", func() {
			It("should return a non-nil result", func() {
				txHash := common.HexToHash("0x28bfa2817bffc06dab8aa8bc2b5524de32bd4296932b25132e86593bc368a8d4")

				trace, err := querier.getTransactionTrace(txHash)
				Expect(err).To(BeNil())
				Expect(trace).ToNot(BeNil())
			})
		})
	})
})
