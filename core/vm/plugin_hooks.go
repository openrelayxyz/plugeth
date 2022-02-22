package vm

func (st *Stack) Len() int {
	return len(st.data)
}
