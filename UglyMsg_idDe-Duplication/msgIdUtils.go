package UglyMsg_idDe_Duplication

import (
	"sort"
)

// MsgSet 定义了一个结构体，用于存储msg_id并提供去重功能。
type MsgSet struct {
	set  map[string]struct{}
	keys []string
}

// NewMsgSet 初始化 MsgSet 结构体。
func NewMsgSet() *MsgSet {
	return &MsgSet{
		set:  make(map[string]struct{}),
		keys: make([]string, 0),
	}
}

// Add 添加一个新的msg_id到集合中。如果集合大小达到100，自动删除最早的一半msg_id。
func (m *MsgSet) Add(msg_id string) bool {
	if _, exists := m.set[msg_id]; exists {
		// msg_id已存在
		return false
	}
	if len(m.set) >= 100 {
		// 如果集合已满，删除最早的一半msg_id。
		m.RemoveHalf()
	}
	m.set[msg_id] = struct{}{}
	m.keys = append(m.keys, msg_id)
	return true
}

// RemoveHalf 删除集合中最早插入的一半msg_id。
func (m *MsgSet) RemoveHalf() {
	sort.Strings(m.keys) // 确保顺序一致
	half := len(m.keys) / 2
	for _, key := range m.keys[:half] {
		delete(m.set, key)
	}
	m.keys = m.keys[half:] // 保留最新的一半msg_id
}

// Exists 检查msg_id是否已经存在于集合中。
func (m *MsgSet) Exists(msg_id string) bool {
	_, exists := m.set[msg_id]
	return exists
}
