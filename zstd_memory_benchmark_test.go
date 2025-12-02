package zstd

import (
	"fmt"
	"testing"
)

// BenchmarkCCtxMemoryUsageByLevel は圧縮レベル別のC側コンテキストのメモリ使用量を測定
func BenchmarkCCtxMemoryUsageByLevel(b *testing.B) {
	levels := []int{1, 3, 6, 9, 12, 15, 19, 22}
	data := generateBenchmarkData(100000) // 100KB

	for _, level := range levels {
		b.Run(fmt.Sprintf("Level%d", level), func(b *testing.B) {
			ctx := NewCtx().(*ctx)

			// 初期状態のメモリサイズを取得
			cctxSizeBefore := ctx.GetSizeofCCtx()

			// 圧縮を実行してコンテキストを初期化
			compressed, err := ctx.CompressLevel(nil, data, level)
			if err != nil {
				b.Fatal(err)
			}

			// 圧縮後のメモリサイズを測定
			cctxSizeAfter := ctx.GetSizeofCCtx()

			// メトリクスの報告
			b.ReportMetric(float64(cctxSizeBefore), "cctx_bytes_before")
			b.ReportMetric(float64(cctxSizeAfter), "cctx_bytes_after")
			b.ReportMetric(float64(cctxSizeAfter-cctxSizeBefore), "cctx_bytes_growth")
			b.ReportMetric(float64(len(compressed)), "compressed_size")
			compressionRatio := float64(len(data)) / float64(len(compressed))
			b.ReportMetric(compressionRatio, "compression_ratio")
		})
	}
}

// BenchmarkDCtxMemoryUsage は解凍コンテキストのC側メモリ使用量を測定
func BenchmarkDCtxMemoryUsage(b *testing.B) {
	data := generateBenchmarkData(100000) // 100KB
	compressed, err := Compress(nil, data)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("Decompress", func(b *testing.B) {
		ctx := NewCtx().(*ctx)

		// 初期状態のメモリサイズを取得
		dctxSizeBefore := ctx.GetSizeofDCtx()

		// 解凍を実行してコンテキストを初期化
		decompressed, err := ctx.Decompress(nil, compressed)
		if err != nil {
			b.Fatal(err)
		}

		// 解凍後のメモリサイズを測定
		dctxSizeAfter := ctx.GetSizeofDCtx()

		// メトリクスの報告
		b.ReportMetric(float64(dctxSizeBefore), "dctx_bytes_before")
		b.ReportMetric(float64(dctxSizeAfter), "dctx_bytes_after")
		b.ReportMetric(float64(dctxSizeAfter-dctxSizeBefore), "dctx_bytes_growth")
		b.ReportMetric(float64(len(decompressed)), "decompressed_size")
	})
}

// BenchmarkMemoryUsageWithMultipleSizes は異なるデータサイズでのメモリ使用量を測定
func BenchmarkMemoryUsageWithMultipleSizes(b *testing.B) {
	sizes := []int{1000, 10000, 100000, 1000000} // 1KB, 10KB, 100KB, 1MB
	level := 6                                   // デフォルトの圧縮レベル

	for _, size := range sizes {
		b.Run(formatSize(size), func(b *testing.B) {
			data := generateBenchmarkData(size)
			ctx := NewCtx().(*ctx)

			// 初期状態のメモリサイズを取得
			cctxSizeBefore := ctx.GetSizeofCCtx()

			// 圧縮を実行
			compressed, err := ctx.CompressLevel(nil, data, level)
			if err != nil {
				b.Fatal(err)
			}

			// 圧縮後のメモリサイズを測定
			cctxSizeAfter := ctx.GetSizeofCCtx()

			// メトリクスの報告
			b.ReportMetric(float64(cctxSizeBefore), "cctx_bytes_before")
			b.ReportMetric(float64(cctxSizeAfter), "cctx_bytes_after")
			b.ReportMetric(float64(cctxSizeAfter-cctxSizeBefore), "cctx_bytes_growth")
			b.ReportMetric(float64(len(compressed)), "compressed_size")
			b.ReportMetric(float64(size), "original_size")
			compressionRatio := float64(size) / float64(len(compressed))
			b.ReportMetric(compressionRatio, "compression_ratio")
		})
	}
}

// BenchmarkMemoryUsagePerformance は圧縮パフォーマンスとメモリ使用量の関係を測定
func BenchmarkMemoryUsagePerformance(b *testing.B) {
	levels := []int{1, 6, 15, 22}
	data := generateBenchmarkData(100000) // 100KB

	for _, level := range levels {
		b.Run(fmt.Sprintf("Level%d", level), func(b *testing.B) {
			b.SetBytes(int64(len(data)))

			// 初回実行でメモリ使用量を測定
			ctx := NewCtx().(*ctx)
			_, err := ctx.CompressLevel(nil, data, level)
			if err != nil {
				b.Fatal(err)
			}
			cctxSizeAfter := ctx.GetSizeofCCtx()

			// メモリメトリクスの報告
			b.ReportMetric(float64(cctxSizeAfter), "cctx_bytes")

			// パフォーマンス測定
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := ctx.CompressLevel(nil, data, level)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// ヘルパー関数

func generateBenchmarkData(size int) []byte {
	data := make([]byte, size)
	// 現実的なデータを模倣するため、パターンを持たせる
	for i := 0; i < size; i++ {
		data[i] = byte((i % 256))
	}
	return data
}

func formatSize(size int) string {
	switch size {
	case 1000:
		return "1KB"
	case 10000:
		return "10KB"
	case 100000:
		return "100KB"
	case 1000000:
		return "1MB"
	default:
		return fmt.Sprintf("%dB", size)
	}
}
