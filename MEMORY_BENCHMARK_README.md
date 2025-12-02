# Zstd Memory Usage Benchmarks

このベンチマークファイルは、zstdライブラリのC側でのメモリ使用量を測定します。

## 含まれるベンチマーク

### 1. BenchmarkCCtxMemoryUsageByLevel
圧縮レベル別（1, 3, 6, 9, 12, 15, 19, 22）のC側圧縮コンテキストのメモリ使用量を測定します。

測定されるメトリクス:
- `cctx_bytes_before`: 圧縮前のコンテキストメモリサイズ
- `cctx_bytes_after`: 圧縮後のコンテキストメモリサイズ
- `cctx_bytes_growth`: メモリの増加量
- `compressed_size`: 圧縮後のデータサイズ
- `compression_ratio`: 圧縮率

### 2. BenchmarkDCtxMemoryUsage
解凍コンテキストのC側メモリ使用量を測定します。

測定されるメトリクス:
- `dctx_bytes_before`: 解凍前のコンテキストメモリサイズ
- `dctx_bytes_after`: 解凍後のコンテキストメモリサイズ
- `dctx_bytes_growth`: メモリの増加量
- `decompressed_size`: 解凍後のデータサイズ

### 3. BenchmarkMemoryUsageWithMultipleSizes
異なるデータサイズ（1KB, 10KB, 100KB, 1MB）でのメモリ使用量を測定します。

### 4. BenchmarkMemoryUsagePerformance
圧縮パフォーマンスとメモリ使用量の関係を測定します（レベル1, 6, 15, 22）。

## 実行方法

すべてのメモリベンチマークを実行:
```bash
cd /path/to/zstd
go test -bench='Benchmark.*Memory' -run=^$ -benchtime=1x
```

特定のベンチマークのみ実行:
```bash
# 圧縮レベル別メモリ使用量
go test -bench=BenchmarkCCtxMemoryUsageByLevel -run=^$ -benchtime=1x

# 解凍メモリ使用量
go test -bench=BenchmarkDCtxMemoryUsage -run=^$ -benchtime=1x

# データサイズ別メモリ使用量
go test -bench=BenchmarkMemoryUsageWithMultipleSizes -run=^$ -benchtime=1x

# パフォーマンス測定
go test -bench=BenchmarkMemoryUsagePerformance -run=^$ -benchtime=1x
```

より詳細な実行（複数回実行して平均を取る）:
```bash
go test -bench='Benchmark.*Memory' -run=^$ -benchtime=10x -count=3
```

## 結果の見方

ベンチマーク結果の例:
```
BenchmarkCCtxMemoryUsageByLevel/Level6-8    1    92125 ns/op    1055984 cctx_bytes_after    5280 cctx_bytes_before    1050704 cctx_bytes_growth    279.0 compressed_size    358.4 compression_ratio
```

これは以下を意味します:
- 圧縮レベル6を使用
- 実行時間: 92,125 ナノ秒
- 圧縮後のコンテキストサイズ: 1,055,984 バイト (約1MB)
- 圧縮前のコンテキストサイズ: 5,280 バイト
- メモリ増加量: 1,050,704 バイト
- 圧縮後のサイズ: 279 バイト
- 圧縮率: 358.4倍

## 重要な発見

圧縮レベルが上がるにつれて、C側のメモリ使用量も増加します:
- Level 1: 約433KB
- Level 6: 約1MB
- Level 15: 約2.7MB
- Level 22: 約2.7MB (Level 15と同じ)

解凍コンテキストは圧縮コンテキストよりもはるかに小さい（約96KB）です。
