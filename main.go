package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

type TextChunk struct {
	text    string
	entropy float64
}

func main() {
	var rootCmd = &cobra.Command{Use: "entropysort"}
	var linesCmd = &cobra.Command{
		Use:   "line",
		Short: "Sort input text by line",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)
			pieces := extractByLine(reader)
			sort.Slice(pieces, func(i, j int) bool {
				return pieces[i].entropy < pieces[j].entropy
			})
			for _, chunk := range pieces {
				fmt.Println(chunk.text)
			}
		},
	}
	var chunksCmd = &cobra.Command{
		Use:   "chunk",
		Short: "Sort input text by byte chunk",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)
			pieces := extractByChunk(128, reader)
			sort.Slice(pieces, func(i, j int) bool {
				return pieces[i].entropy < pieces[j].entropy
			})
			for _, chunk := range pieces {
				fmt.Println(chunk.text)
			}
		},
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(linesCmd)
	rootCmd.AddCommand(chunksCmd)
	rootCmd.Execute()
}

func entropy(s string) float64 {
	var characterMap = make(map[rune]int)
	for _, r := range s {
		characterMap[r]++
	}
	var prob float64
	var result float64
	for _, c := range characterMap {
		prob = float64(c) / float64(len(s))
		result -= prob * math.Log2(prob)
	}
	return result
}

func extractByLine(reader io.Reader) []TextChunk {
	scanner := bufio.NewScanner(reader)
	var pieces []TextChunk
	for scanner.Scan() {
		value := scanner.Text()
		pieces = append(pieces, TextChunk{text: value, entropy: entropy(strings.TrimSpace(value))})
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return pieces
}

func extractByChunk(size int, reader io.Reader) []TextChunk {
	p := make([]byte, size)
	var pieces []TextChunk
	for {
		_, err := reader.Read(p)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		value := string(p)
		pieces = append(pieces, TextChunk{text: value, entropy: entropy(strings.TrimSpace(value))})
	}
	return pieces
}
