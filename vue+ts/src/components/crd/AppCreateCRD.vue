<template>
  <div class="bookapp-create-page">
    <div class="page-container">
      <h2 class="title">
        <i class="fa fa-book icon"></i> 创建 Bookapp
      </h2>

      <!-- 只读基本信息 -->
      <div class="readonly-info">
        <p><strong>apiVersion：</strong> webapp.y2505.com/v1</p>
        <p><strong>kind：</strong> Bookapp</p>
      </div>

      <div class="main-content">
        <!-- 左边表单 -->
        <div class="form-content">
          <!-- 基本信息 -->
          <div class="form-section">
            <h3>
              <i class="fa fa-info-circle section-icon blue"></i> 基本信息
            </h3>
            <div class="form-item">
              <label>名称：</label>
              <input v-model="bookapp.name" type="text" placeholder="请输入 Bookapp 名称" />
            </div>
            <div class="form-item">
              <label>命名空间：</label>
              <select v-model="bookapp.namespace">
                <option value="">请选择命名空间</option>
                <option v-for="ns in namespaces" :key="ns" :value="ns">{{ ns }}</option>
              </select>
            </div>
          </div>

          <!-- Labels -->
          <div class="form-section">
            <h3>
              <i class="fa fa-tag section-icon orange"></i> Labels
            </h3>
            <p class="labels-hint">⚠️ Labels 将自动生成到 metadata.labels 中</p>
            <div v-for="(label, idx) in labelsArray" :key="'label-' + idx" class="form-item">
              <label>Key：</label>
              <input v-model="label.key" type="text" placeholder="如 app" @input="updateLabels" />
              <label>Value：</label>
              <input v-model="label.value" type="text" placeholder="如 bookapp" @input="updateLabels" />
              <button class="btn btn-danger btn-small" @click="removeLabel(idx)">
                <i class="fa fa-trash"></i>
              </button>
            </div>
            <button class="btn btn-add" @click="addLabel">
              <i class="fa fa-plus"></i> 添加 Label
            </button>
          </div>

          <!-- Spec 配置 -->
          <div class="form-section">
            <h3>
              <i class="fa fa-cogs section-icon purple"></i> Spec 配置
            </h3>
            <div class="form-item">
              <label>副本数：</label>
              <input type="number" v-model="bookapp.spec.size" placeholder="请输入副本数" />
            </div>
            <div class="form-item">
              <label>镜像：</label>
              <input type="text" v-model="bookapp.spec.image" placeholder="如 nginx:latest" />
            </div>
            <div class="form-item">
              <label>镜像拉取策略：</label>
              <select v-model="bookapp.spec.imagePullPolicy">
               <option value="IfNotPresent">IfNotPresent</option>
              <option value="Always">Always</option>
              <option value="Never">Never</option>
             </select>
            </div>
            <div class="form-item">
              <label>端口名称：</label>
              <input type="text" v-model="bookapp.spec.port.name" placeholder="如 web" />
            </div>
            <div class="form-item">
              <label>协议：</label>
              <select v-model="bookapp.spec.port.protocol">
                <option value="TCP">TCP</option>
                <option value="UDP">UDP</option>
              </select>
            </div>
            <div class="form-item">
              <label>端口：</label>
              <input type="number" v-model="bookapp.spec.port.port" placeholder="端口号" />
            </div>
            <div class="form-item">
              <label>目标端口：</label>
              <input type="number" v-model="bookapp.spec.port.targetPort" placeholder="目标端口号" />
            </div>
            <div class="form-item">
              <label>域名：</label>
              <input type="text" v-model="bookapp.spec.serverName" placeholder="如 www.example.com" />
            </div>
          </div>
        </div>

        <!-- YAML 预览 -->
        <div class="yaml-preview">
          <h3><i class="fa fa-file-code-o"></i> YAML 预览</h3>
          <button class="copy-btn" @click="copyYaml">复制</button>
          <pre>{{ yamlOutput }}</pre>
        </div>
      </div>

      <!-- 提交按钮 -->
      <div class="form-actions">
        <button class="btn btn-primary" @click="createBookapp">
          <i class="fa fa-check"></i> 创建 Bookapp
        </button>
        <button class="btn btn-secondary" @click="resetForm">
          <i class="fa fa-refresh"></i> 重置
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref, computed } from "vue"
import { useK8sData } from "@/js/useK8sData"
import yaml from "js-yaml"
//import axios from "axios"  // 如果还没导入
import axios from "@/main" 

export default defineComponent({
  name: "AppBookappCreate",
  setup() {
    const { namespaces } = useK8sData()

    const bookapp = ref({
      name: "",
      namespace: "",
      labels: {},
      spec: {
        size: 1,
        image: "",
        serverName: "",
        port: { name: "", protocol: "TCP", port: 80, targetPort: 80 },
        imagePullPolicy: "IfNotPresent", // ✅ 默认值
      }
    })


    const labelsArray = ref([{ key: "", value: "" }])

    const updateLabels = () => {
      const obj = {}
      labelsArray.value.forEach(l => {
        if (l.key) obj[l.key] = l.value
      })
      bookapp.value.labels = obj
    }

    const addLabel = () => labelsArray.value.push({ key: "", value: "" })
    const removeLabel = idx => { labelsArray.value.splice(idx, 1); updateLabels() }
    const resetForm = () => {
      bookapp.value = {
        name: "",
        namespace: "",
        labels: {},
        spec: {
          size: 1,
          image: "",
          serverName: "",
          port: { name: "", protocol: "TCP", port: 80, targetPort: 80 },
          imagePullPolicy: "IfNotPresent" // ✅ 添加默认值
        }
      }
      labelsArray.value = [{ key: "", value: "" }]
    }

    const yamlOutput = computed(() => yaml.dump(bookapp.value, { noRefs: true }))

    const copyYaml = async () => {
      try { await navigator.clipboard.writeText(yamlOutput.value); alert("已复制到剪贴板") }
      catch { alert("复制失败，请手动选择复制") }
    }

 const createBookapp = async () => {
  if (!bookapp.value.name || !bookapp.value.namespace) {
    return alert("❌ 请填写名称和命名空间")
  }

  try {
    // 构造请求 payload
    const payload = {
      name: bookapp.value.name,
      namespace: bookapp.value.namespace,
      labels: bookapp.value.labels,
      spec: {
        size: Number(bookapp.value.spec.size),
        image: bookapp.value.spec.image,
        imagePullPolicy: bookapp.value.spec.imagePullPolicy || "IfNotPresent",
        port: {
          name: bookapp.value.spec.port.name,
          protocol: bookapp.value.spec.port.protocol,
          port: Number(bookapp.value.spec.port.port),
          targetPort: Number(bookapp.value.spec.port.targetPort)
        },
        serverName: bookapp.value.spec.serverName
      }
    }

    const res = await axios.post("http://192.168.216.50:8090/bookapp/", payload)

    if (res.data && res.data.message) {
      alert("✅ " + res.data.message)
    } else {
      alert("✅ 创建成功")
    }

    resetForm()
  } catch (err) {
    console.error("创建 Bookapp 失败:", err)
    alert("❌ 创建失败，请检查后台日志")
  }
}


    return {
      namespaces, bookapp, labelsArray,
      addLabel, removeLabel, updateLabels,
      resetForm, yamlOutput, copyYaml, createBookapp
    }
  }
})
</script>

<style scoped>
.page-container { background: #fff; border-radius: 12px; padding: 20px; box-shadow: 0 4px 10px rgba(0,0,0,0.1); }
.title>.icon { color: #42a5f5 }
.readonly-info { background: #f5f5f5; padding: 8px 12px; border-radius: 6px; margin-bottom: 15px; font-size: 14px }
.main-content { display: flex; gap: 20px }
.form-content { flex: 1 }
.yaml-preview {
  flex: 1;
  background: #fff;
  color: #333;
  padding: 16px;
  border-radius: 8px;
  font-size: 15px;
  overflow-x: auto;

  /* Pod 风格的渐变边框 */
  border: 2px solid transparent;
  background-clip: padding-box;
  position: relative;
}
.yaml-preview::before {
  content: "";
  position: absolute;
  inset: 0;
  border-radius: 8px;
  padding: 2px;
  background: linear-gradient(45deg, #42a5f5, #66bb6a);
  -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
  mask-composite: exclude;
}

.copy-btn { position: absolute; top: 8px; right: 12px; background: #42a5f5; color: #fff; border: none; padding: 4px 10px; border-radius: 6px; cursor: pointer; font-size: 12px }
.form-section { margin-bottom: 20px; padding: 15px; border: 1px solid #eee; border-radius: 10px }
.section-icon.blue { color: #42a5f5 }
.section-icon.orange { color: #ff9800 }
.section-icon.purple { color: #9c27b0 }
.form-item { display: flex; align-items: center; margin: 10px 0; gap: 8px }
.form-item label { width: 80px; font-weight: 500; color: #555 }
.form-item input, .form-item select { flex: 1; padding: 6px 10px; border: 1px solid #ddd; border-radius: 6px }
.form-actions { margin-top: 20px; display: flex; gap: 12px }
.btn { padding: 6px 12px; border: none; border-radius: 8px; cursor: pointer; font-size: 14px }
.btn-primary { background: #42a5f5; color: #fff }
.btn-secondary { background: #ccc; color: #333 }
.btn-add { background: #3178c6; color: #fff; margin-top: 10px }
.btn-danger { background: #f66d6d; color: #fff }
.btn-small { font-size: 12px; padding: 4px 8px }
.labels-hint { font-size: 12px; color: #888; margin-bottom: 6px }
</style>
