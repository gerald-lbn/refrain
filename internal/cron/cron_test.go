package scheduler_test

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	scheduler "github.com/gerald-lbn/refrain/internal/cron"
	"github.com/gerald-lbn/refrain/internal/domain"
)

var _ = Describe("Scheduler", func() {
	var (
		sched  domain.Scheduler
		logger *slog.Logger
	)

	BeforeEach(func() {
		logger = slog.New(slog.NewTextHandler(GinkgoWriter, &slog.HandlerOptions{Level: slog.LevelDebug}))
		sched = scheduler.New(logger)
	})

	Describe("New", func() {
		It("should create a scheduler", func() {
			Expect(sched).NotTo(BeNil())
		})
	})

	Describe("Start", func() {
		It("should not panic", func() {
			Expect(func() { sched.Start(context.Background()) }).NotTo(Panic())
		})
	})

	Describe("AddFunc", func() {
		It("should execute the function at the given interval", func() {
			var count atomic.Int32

			sched.AddFunc(50*time.Millisecond, func() {
				count.Add(1)
			})

			Eventually(func() int32 {
				return count.Load()
			}).WithTimeout(500 * time.Millisecond).WithPolling(10 * time.Millisecond).Should(BeNumerically(">=", 2))

			sched.Stop()
		})

		It("should support multiple scheduled functions", func() {
			var countA, countB atomic.Int32

			sched.AddFunc(50*time.Millisecond, func() {
				countA.Add(1)
			})
			sched.AddFunc(50*time.Millisecond, func() {
				countB.Add(1)
			})

			Eventually(func() int32 {
				return countA.Load()
			}).WithTimeout(500 * time.Millisecond).WithPolling(10 * time.Millisecond).Should(BeNumerically(">=", 1))

			Eventually(func() int32 {
				return countB.Load()
			}).WithTimeout(500 * time.Millisecond).WithPolling(10 * time.Millisecond).Should(BeNumerically(">=", 1))

			sched.Stop()
		})
	})

	Describe("Stop", func() {
		It("should stop all scheduled functions", func() {
			var count atomic.Int32

			sched.AddFunc(50*time.Millisecond, func() {
				count.Add(1)
			})

			Eventually(func() int32 {
				return count.Load()
			}).WithTimeout(500 * time.Millisecond).WithPolling(10 * time.Millisecond).Should(BeNumerically(">=", 1))

			sched.Stop()
			snapshot := count.Load()

			time.Sleep(200 * time.Millisecond)
			Expect(count.Load()).To(Equal(snapshot))
		})
	})
})
