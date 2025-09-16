// src/js/useK8sData.js
import { ref, onMounted } from "vue"
import axios from "axios"

export function useK8sData() {
  const namespaces = ref([])
  const nodes = ref([])
  const loading = ref(false)
  const error = ref(null)

  const fetchNamespaces = async () => {
    try {
      const res = await axios.get("http://192.168.216.50:8090/getnamespacelist")
      namespaces.value = res.data.data || []
    } catch (err) {
      console.error("获取 namespace 失败:", err)
      error.value = err
    }
  }

  const fetchNodes = async () => {
    try {
      const res = await axios.get("http://192.168.216.50:8090/nodemetrics")
      nodes.value = res.data.data.map(n => n.name)
    } catch (err) {
      console.error("获取 node 失败:", err)
      error.value = err
    }
  }

  const fetchAll = async () => {
    loading.value = true
    await Promise.all([fetchNamespaces(), fetchNodes()])
    loading.value = false
  }

  onMounted(fetchAll)

  return { namespaces, nodes, loading, error, fetchAll }
}
