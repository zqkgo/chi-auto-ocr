package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitLines(t *testing.T) {
	txt := `蛋白质是生命的物质基础，是有机大分子，是构成细胞的基本有机物，是生命活动的主要承担者。没有蛋白质就没有生命。氨基酸是蛋白质的基本组成单位。它是与生命及与各种形式的生命活动紧密联系在一起的物质。机体中的每一个细胞和所有重要组成部分都有蛋白质参与。蛋白质占人体重量的16%~20%，即一个60kg重的成年人其体内约有蛋白质9.6~12kg。人体内蛋白质的种类很多，性质、功能各异，但都是由20多种氨基酸（Amino acid）按不同比例组合而成的，并在体内不断进行代谢与更新。`
	txt = splitLines(txt)
	lines := strings.Split(txt, "\n\n")
	require.Equal(t, 7, len(lines))

	txt = `UDP uses a simple connectionless communication model with a minimum of protocol mechanisms. UDP provides checksums for data integrity, and port numbers for addressing different functions at the source and destination of the datagram. It has no handshaking dialogues, and thus exposes the user's program to any unreliability of the underlying network; there is no guarantee of delivery, ordering, or duplicate protection. If error-correction facilities are needed at the network interface level, an application may instead use Transmission Control Protocol (TCP) or Stream Control Transmission Protocol (SCTP) which are designed for this purpose.`
	txt = splitLines(txt)
	lines = strings.Split(txt, "\n\n")
	require.Equal(t, 4, len(lines))

	txt = ""
	txt = splitLines(txt)
	require.Empty(t, txt)

	txt = "床前明月光"
	txt = splitLines(txt)
	require.NotEmpty(t, txt)

	txt = "低头思故乡。"
	txt2 := splitLines(txt)
	require.Equal(t, txt, txt2)
}
