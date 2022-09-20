package rbtree

import "fmt"

// PrintTree 打印树
func PrintTree(t *RBTree) {
	const (
		nilStr = "nil"
		indent = 2
	)
	levelNode := make(map[int][]*node)
	levelNode[0] = []*node{t.root}
	for level := 0; ; level++ {
		var nodes = levelNode[level]
		var next []*node
		for _, n := range nodes {
			if n != nil {
				next = append(next, n.left, n.right)
			} else {
				next = append(next, nil, nil)
			}
		}
		var exit = true
		for _, v := range next {
			if v != nil {
				exit = false
				break
			}
		}
		if exit {
			break
		}
		levelNode[level+1] = next
	}
	depth := len(levelNode)
	for j := 0; j < depth; j++ {
		w := indent << (depth - 1 - j)
		if j > 0 {
			for i := 0; i < 1<<(j-1); i++ {
				if levelNode[j][i*2] == nil {
					fmt.Printf("%*c", w*4, ' ')
				} else {
					fmt.Printf("%*c", w, ' ') // w
					if w < 3 {
						leftW := 3
						if w == 1 {
							fmt.Printf("| ")
						} else {
							fmt.Printf("/ \\")
						}
						leftW -= 3 / w
						n := w - 3%w + leftW
						fmt.Printf("%*c", n, ' ')
					} else {
						fmt.Printf("%c", ' ')
						for k := 0; k < w-3; k++ {
							fmt.Printf("_")
						}
						fmt.Printf("/ \\")
						for k := 0; k < w-3; k++ {
							fmt.Printf("_")
						}
						fmt.Printf("%*c", w+2, ' ')
					}
				}
			}

			fmt.Printf("\n")
			for i := 0; i < 1<<(j-1); i++ {
				if levelNode[j][i*2] == nil {
					fmt.Printf("%*c", w*4, ' ')
				} else {
					if w < 3 {
						fmt.Printf("%*c%*c%*c", w, '/', w*2, '\\', w, ' ')
					} else {
						fmt.Printf("%*c%*c%*c", w+1, '/', w*2-2, '\\', w+1, ' ')
					}
				}
			}
			fmt.Printf("\n")
		}
		for i := 0; i < 1<<j; i++ {
			n := levelNode[j][i]
			if n == nil {
				fmt.Printf("%*c", w*2, ' ')
				continue
			}
			key := fmt.Sprintf("%v", n.item)
			if n.item == nil {
				key = nilStr
			}
			shiftLeft := (len(key) + 1) / 2
			if w < 3 {
				if i%2 == 0 || len(key) > 2 {
					shiftLeft = (len(key))/2 + 1
				} else {
					shiftLeft = (len(key) + 1) / 2
				}
			}
			if shiftLeft > w {
				shiftLeft = w
			}
			if w > shiftLeft {
				fmt.Printf("%*c", w-shiftLeft, ' ') // (key)
			}
			if n.color == RED {
				fmt.Printf("%c[1;41;37m%v%c[0m", 0x1B, key, 0x1B)
			} else {
				fmt.Printf("%c[1;40;37m%v%c[0m", 0x1B, key, 0x1B)
			}
			fmt.Printf("%*c", w-(len(key)-shiftLeft), ' ')
		}
		fmt.Printf("\n")
	}
}
