# Introduction

The main goal of this project is to explore how compression algorithms work and how to build them in go.

The name `janus-zip` comes from the Roman/Etruscan god mainly associated with opening and closing, the start and the beginning. As zipping a file is like backing it together in a compact box (closing) and then after transferring it to the correct place opening this compact box again I found the name fitting (https://en.wikipedia.org/wiki/Janus).

This repository uses basic Huffman coding to compress files down, this tries to shrink down the amount of bits used for characters that occur more often, how it exactly works can be found here: https://en.wikipedia.org/wiki/Huffman_coding