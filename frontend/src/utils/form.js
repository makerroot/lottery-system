/**
 * 表单数据工具函数
 * 用于自动去除表单数据中的前后空格
 */

/**
 * 去除对象中所有字符串字段的前后空格
 * @param {Object} obj - 表单数据对象
 * @returns {Object} 处理后的对象
 */
export function trimObject(obj) {
  if (!obj || typeof obj !== 'object') {
    return obj
  }

  const trimmed = {}

  for (const key in obj) {
    if (obj.hasOwnProperty(key)) {
      const value = obj[key]

      if (typeof value === 'string') {
        // 字符串类型：去除前后空格
        trimmed[key] = value.trim()
      } else if (Array.isArray(value)) {
        // 数组类型：递归处理每个元素
        trimmed[key] = value.map(item =>
          typeof item === 'string' ? item.trim() : trimObject(item)
        )
      } else if (value !== null && typeof value === 'object') {
        // 对象类型：递归处理
        trimmed[key] = trimObject(value)
      } else {
        // 其他类型（数字、布尔等）：保持不变
        trimmed[key] = value
      }
    }
  }

  return trimmed
}

/**
 * 去除数组中所有字符串元素的前后空格
 * @param {Array} arr - 字符串数组
 * @returns {Array} 处理后的数组
 */
export function trimArray(arr) {
  if (!Array.isArray(arr)) {
    return arr
  }

  return arr.map(item =>
    typeof item === 'string' ? item.trim() : item
  )
}

/**
 * 为Ant Design Vue表单组件添加自动trim功能
 * 在v-model时使用此函数
 * @param {String} value - 输入值
 * @returns {String} trim后的值
 */
export function autoTrim(value) {
  if (typeof value !== 'string') {
    return value
  }
  return value.trim()
}
