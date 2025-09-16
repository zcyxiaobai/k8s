<template>
  <div class="deployment-create-page">
    <div class="page-container">
      <h2 class="title">
        <i class="fa fa-cubes icon"></i> 创建 Deployment
      </h2>

      <!-- 只读的基本固定字段 -->
      <div class="readonly-info">
        <p><strong>apiVersion：</strong> apps/v1</p>
        <p><strong>kind：</strong> Deployment</p>
      </div>

      <div class="main-content">
        <!-- 左边表单 -->
        <div class="form-content">
          <!-- Deployment 基本信息 -->
          <div class="form-section">
            <h3>
              <i class="fa fa-info-circle section-icon blue"></i> 基本信息
            </h3>
            <div class="form-item">
              <label>名称：</label>
              <input v-model="deployment.name" type="text" placeholder="请输入 Deployment 名称" />
            </div>
            <div class="form-item">
              <label>命名空间：</label>
              <select v-model="deployment.namespace">
                <option value="">请选择命名空间</option>
                <option v-for="ns in namespaces" :key="ns" :value="ns">{{ ns }}</option>
              </select>
            </div>
            <div class="form-item">
              <label>副本数：</label>
              <input v-model="deployment.replicas" type="number" min="1" placeholder="请输入副本数" />
            </div>
          </div>

          <!-- Deployment 自身 Labels -->
          <div class="form-section">
            <h3>
              <i class="fa fa-tag section-icon orange"></i> Deployment Labels (metadata.labels)
            </h3>
            <div v-for="(label, idx) in deployment.metaLabels" :key="'meta-' + idx" class="form-item">
              <label>Key：</label>
              <input v-model="label.key" type="text" placeholder="如 env" />
              <label>Value：</label>
              <input v-model="label.value" type="text" placeholder="如 prod" />
              <button class="btn btn-danger btn-small" @click="removeMetaLabel(idx)">
                <i class="fa fa-trash"></i>
              </button>
            </div>
            <button class="btn btn-add" @click="addMetaLabel">
              <i class="fa fa-plus"></i> 添加 Metadata Label
            </button>
          </div>

          <!-- Pod Selector & Template Labels -->
          <div class="form-section">
            <h3>
              <i class="fa fa-tags section-icon green"></i> Pod Labels (selector.matchLabels & template.metadata.labels)
            </h3>
            <p class="labels-hint">⚠️ 此处的 Labels 将自动同步到 selector.matchLabels 与 template.metadata.labels</p>
            <div v-for="(label, idx) in deployment.labels" :key="'pod-' + idx" class="form-item">
              <label>Key：</label>
              <input v-model="label.key" type="text" placeholder="如 app" />
              <label>Value：</label>
              <input v-model="label.value" type="text" placeholder="如 nginx" />
              <button class="btn btn-danger btn-small" @click="removeLabel(idx)">
                <i class="fa fa-trash"></i>
              </button>
            </div>
            <button class="btn btn-add" @click="addLabel">
              <i class="fa fa-plus"></i> 添加 Pod Label
            </button>
          </div>

          <!-- 容器配置 -->
          <div class="form-section">
            <h3>
              <i class="fa fa-box section-icon purple"></i> 容器配置
            </h3>
            <div v-for="(container, idx) in deployment.containers" :key="idx" class="container-card">
              <h4><i class="fa fa-cube container-icon"></i> 容器 {{ idx + 1 }}</h4>
              <div class="form-item">
                <label>名称：</label>
                <input v-model="container.name" type="text" placeholder="容器名称" />
              </div>
              <div class="form-item">
                <label>镜像：</label>
                <input v-model="container.image" type="text" placeholder="如 nginx:latest" />
              </div>
              <div class="form-item">
                <label>镜像拉取策略：</label>
                <select v-model="container.imagePullPolicy">
                  <option value="">请选择策略</option>
                  <option value="IfNotPresent">IfNotPresent</option>
                  <option value="Always">Always</option>
                  <option value="Never">Never</option>
                </select>
              </div>
              <div class="form-item">
                <label>端口：</label>
                <input v-model="container.port" type="number" placeholder="容器端口 (可选)" />
              </div>
              <button class="btn btn-danger btn-small" @click="removeContainer(idx)">
                <i class="fa fa-trash"></i> 移除容器
              </button>
            </div>
            <button class="btn btn-add" @click="addContainer">
              <i class="fa fa-plus"></i> 添加容器
            </button>
          </div>
        </div>

        <!-- 右边 YAML 实时预览 -->
        <div class="yaml-preview">
          <h3><i class="fa fa-file-code-o"></i> YAML 预览</h3>
          <button class="copy-btn" @click="copyYaml">复制</button>
          <pre>{{ yamlOutput }}</pre>
        </div>
      </div>

      <!-- 提交按钮 -->
      <div class="form-actions">
        <button class="btn btn-primary" @click="createDeployment">
          <i class="fa fa-check"></i> 创建 Deployment
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
// import axios from "axios"
import axios from "@/main" 
import { useK8sData } from "@/js/useK8sData"
import yaml from "js-yaml"

export default defineComponent({
  name: "AppDeploymentCreate",
  setup() {
    const { namespaces } = useK8sData()

    const deployment = ref({
      name: "",
      namespace: "",
      replicas: "",
      metaLabels: [{ key: "env", value: "prod" }], // metadata.labels
      labels: [{ key: "app", value: "nginx" }],    // selector + template
      containers: [
        { name: "", image: "", imagePullPolicy: "", port: "" }
      ]
    })

    const metaLabelsObj = computed(() => {
      const obj = {}
      deployment.value.metaLabels.forEach(l => {
        if (l.key && l.value) obj[l.key] = l.value
      })
      return obj
    })

    const podLabelsObj = computed(() => {
      const obj = {}
      deployment.value.labels.forEach(l => {
        if (l.key && l.value) obj[l.key] = l.value
      })
      return obj
    })

    // YAML 实时生成，禁止引用 &ref
    const yamlOutput = computed(() => {
      const obj = {
        apiVersion: "apps/v1",
        kind: "Deployment",
        metadata: {
          name: deployment.value.name || "示例-deployment",
          namespace: deployment.value.namespace || "default",
          labels: metaLabelsObj.value
        },
        spec: {
          replicas: deployment.value.replicas || 1,
          selector: {
            matchLabels: podLabelsObj.value
          },
          template: {
            metadata: {
              labels: podLabelsObj.value
            },
            spec: {
              containers: deployment.value.containers.map(c => ({
                name: c.name || "container-demo",
                image: c.image || "nginx:latest",
                imagePullPolicy: c.imagePullPolicy || "IfNotPresent",
                ports: c.port ? [{ containerPort: Number(c.port) }] : []
              }))
            }
          }
        }
      }
      return yaml.dump(obj, { noRefs: true })
    })

    const addContainer = () => deployment.value.containers.push({ name: "", image: "", imagePullPolicy: "", port: "" })
    const removeContainer = idx => deployment.value.containers.splice(idx, 1)
    const addLabel = () => deployment.value.labels.push({ key: "", value: "" })
    const removeLabel = idx => deployment.value.labels.splice(idx, 1)
    const addMetaLabel = () => deployment.value.metaLabels.push({ key: "", value: "" })
    const removeMetaLabel = idx => deployment.value.metaLabels.splice(idx, 1)

    const resetForm = () => {
      deployment.value = {
        name: "",
        namespace: "",
        replicas: "",
        metaLabels: [{ key: "env", value: "prod" }],
        labels: [{ key: "app", value: "nginx" }],
        containers: [{ name: "", image: "", imagePullPolicy: "", port: "" }]
      }
    }

 const createDeployment = async () => {
  try {
    // 构造 payload，确保数字类型正确
    const payload = {
      name: deployment.value.name,
      namespace: deployment.value.namespace,
      replicas: Number(deployment.value.replicas) || 1,
      metaLabels: {},
      labels: {},
      containers: deployment.value.containers.map(c => ({
        name: c.name,
        image: c.image,
        imagePullPolicy: c.imagePullPolicy || "IfNotPresent",
        port: Number(c.port) || 0
      }))
    }

    // metadata.labels
    deployment.value.metaLabels.forEach(l => {
      if (l.key && l.value) payload.metaLabels[l.key] = l.value
    })

    // pod labels
    deployment.value.labels.forEach(l => {
      if (l.key && l.value) payload.labels[l.key] = l.value
    })

    console.log("发送数据:", payload)

    // 请求后端接口
    const res = await axios.post("http://192.168.216.50:8090/dep/CreateDev", payload)

    // 提示创建结果
    alert(res.data.message || "创建成功")
  } catch (err) {
    console.error(err)
    alert("创建失败，请检查控制台")
  }
}



    const copyYaml = async () => {
      try {
        await navigator.clipboard.writeText(yamlOutput.value)
        alert("已复制到剪贴板")
      } catch {
        alert("复制失败，请手动选择复制")
      }
    }

    return {
      namespaces,
      deployment,
      yamlOutput,
      addContainer,
      removeContainer,
      addLabel,
      removeLabel,
      addMetaLabel,
      removeMetaLabel,
      resetForm,
      createDeployment,
      copyYaml
    }
  }
})
</script>

<style scoped>
.title>.icon { color: #ff9800 }
.page-container { background: #fff; border-radius: 12px; padding: 20px; box-shadow: 0 4px 10px rgba(0,0,0,0.1); }
.readonly-info { background: #f5f5f5; padding: 8px 12px; border-radius: 6px; margin-bottom: 15px; font-size: 14px; }
.main-content { display: flex; gap: 20px; }
.form-content { flex: 1 }
.yaml-preview {
  flex: 1;
  background: #fff;
  color: #333;
  padding: 16px;
  border-radius: 8px;
  font-size: 15px;
  overflow-x: auto;
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
.copy-btn { position: absolute; top: 8px; right: 12px; background: #42a5f5; color: #fff; border: none; padding: 4px 10px; border-radius: 6px; cursor: pointer; font-size: 12px; }
.form-section { margin-bottom: 20px; padding: 15px; border: 1px solid #eee; border-radius: 10px; }
.section-icon.blue { color: #42a5f5 }
.section-icon.green { color: #66bb6a }
.section-icon.orange { color: #ff9800 }
.section-icon.purple { color: #9c27b0 }
.container-icon { color: #6a1b9a; margin-right: 6px; }
.form-item { display: flex; align-items: center; margin: 10px 0; gap: 8px; }
.form-item label { width: 60px; font-weight: 500; color: #555; }
.form-item input, .form-item select { flex: 1; padding: 6px 10px; border: 1px solid #ddd; border-radius: 6px; }
.container-card { background: #fafafa; border: 1px solid #ddd; border-radius: 10px; padding: 12px; margin-bottom: 12px; }
.form-actions { margin-top: 20px; display: flex; gap: 12px }
.btn { padding: 6px 12px; border: none; border-radius: 8px; cursor: pointer; font-size: 14px }
.btn-primary { background: #4CAF50; color: #fff }
.btn-secondary { background: #ccc; color: #333 }
.btn-add { background: #3178c6; color: #fff; margin-top: 10px }
.btn-danger { background: #f66d6d; color: #fff }
.btn-small { font-size: 12px; padding: 4px 8px }
.labels-hint { font-size: 12px; color: #888; margin-bottom: 6px; }
</style>
