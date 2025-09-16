<template>
  <div class="pod-update-page">
    <div class="page-container">
      <h2 class="title">
        <i class="fa fa-edit icon"></i> 修改 Pod - {{ podDataCopy?.name }}
      </h2>

      <!-- 只读基本信息 -->
      <div class="readonly-info">
        <p><strong>apiVersion：</strong> v1</p>
        <p><strong>kind：</strong> Pod</p>
        <p><strong>命名空间：</strong> {{ podDataCopy?.namespace }}</p>
      </div>

      <div class="main-content">
        <!-- 左边表单 -->
        <div class="form-content">
          <!-- Labels 可编辑 -->
          <div class="form-section">
            <h3>
              <i class="fa fa-tag section-icon orange"></i> Pod Labels
            </h3>
            <p class="labels-hint">⚠️ 可编辑，修改 metadata.labels</p>
            <div
              v-for="(label, idx) in podLabelsArray"
              :key="'label-' + idx"
              class="form-item"
            >
              <label>Key：</label>
              <input
                v-model="label.key"
                type="text"
                placeholder="如 app"
                @input="syncLabels"
              />
              <label>Value：</label>
              <input
                v-model="label.value"
                type="text"
                placeholder="如 nginx"
                @input="syncLabels"
              />
              <button class="btn btn-danger btn-small" @click="removeLabel(idx)">
                <i class="fa fa-trash"></i>
              </button>
            </div>
            <button class="btn btn-add" @click="addLabel">
              <i class="fa fa-plus"></i> 添加 Label
            </button>
          </div>

          <!-- 容器信息 只读 -->
          <div class="form-section">
            <h3>
              <i class="fa fa-box section-icon purple"></i> 容器信息 (只读)
            </h3>
            <div
              v-for="(container, idx) in podDataCopy?.containers || []"
              :key="'c-' + idx"
              class="container-card"
            >
              <h4>
                <i class="fa fa-cube container-icon"></i> 容器 {{ idx + 1 }}
              </h4>
              <div class="form-item">
                <label>名称：</label>
                <input type="text" :value="container.name" readonly />
              </div>
              <div class="form-item">
                <label>镜像：</label>
                <input type="text" :value="container.image" readonly />
              </div>
              <div class="form-item">
                <label>镜像拉取策略：</label>
                <input type="text" :value="container.imagePullPolicy" readonly />
              </div>
              <div class="form-item">
                <label>端口：</label>
                <input
                  type="text"
                  :value="(container.ports && container.ports.join(',')) || '-'"
                  readonly
                />
              </div>
            </div>
          </div>

          <!-- 其他字段只读 -->
          <div class="form-section">
            <h3>
              <i class="fa fa-info-circle section-icon blue"></i> 其他信息 (只读)
            </h3>
            <div
              v-for="(value, key) in podDataCopy"
              :key="key"
              class="form-item"
            >
              <template v-if="key !== 'labels' && key !== 'containers' && key !== 'name' && key !== 'namespace'">
                <label>{{ key }}：</label>
                <input type="text" :value="value" readonly />
              </template>
            </div>
          </div>
        </div>

        <!-- 右边 YAML 预览 -->
        <div class="yaml-preview">
          <h3><i class="fa fa-file-code-o"></i> Pod 数据预览 (YAML)</h3>
          <button class="copy-btn" @click="copyYaml">复制</button>
          <pre>{{ yamlContent }}</pre>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="form-actions">
        <button class="btn btn-primary" @click="updatePod">
          <i class="fa fa-check"></i> 保存修改
        </button>
        <button class="btn btn-secondary" @click="goBack">
          <i class="fa fa-refresh"></i> 取消
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref, onMounted, reactive, computed } from "vue";
// import axios from "axios";
import axios from "@/main" 
import { useRouter } from "vue-router";
import yaml from "js-yaml";

export default defineComponent({
  name: "UpdatePod",
  props: {
    podData: { type: Object, default: null }
  },
  setup(props) {
    const router = useRouter();
    const podLabelsArray = ref([{ key: "", value: "" }]);

    // 拷贝 podData 用于实时更新 JSON
    const podDataCopy = reactive(JSON.parse(JSON.stringify(props.podData || {})));

    // 初始化 labels
    onMounted(() => {
      if (props.podData?.labels) {
        podLabelsArray.value = Object.entries(props.podData.labels).map(
          ([k, v]) => ({ key: k, value: v })
        );
      }
    });

    // 同步左侧 labels 到 podDataCopy
    const syncLabels = () => {
      const labelsObj = {};
      podLabelsArray.value.forEach(l => {
        if (l.key) labelsObj[l.key] = l.value;
      });
      podDataCopy.labels = labelsObj;
    };

    // 添加 label
    const addLabel = () => {
      podLabelsArray.value.push({ key: "", value: "" });
    };

    // 删除 label
    const removeLabel = (idx) => {
      podLabelsArray.value.splice(idx, 1);
      syncLabels();
    };

    // 更新 pod
   // 更新 pod
const updatePod = async () => {
  // 组装 labels 对象
  const labelsObj = {};
  podLabelsArray.value.forEach((l) => {
    if (l.key) labelsObj[l.key] = l.value;
  });

  try {
    const res = await axios.post("http://192.168.216.50:8090/pod/Update", {
      name: podDataCopy.name,        // Pod 名称
      namespace: podDataCopy.namespace, // 命名空间
      labels: labelsObj              // 新的 labels
    });

    alert(res.data.message || "修改成功");
    router.back();
  } catch (err) {
    console.error(err);
    alert("修改失败，请检查控制台信息");
  }
};

    // 返回
    const goBack = () => router.back();

    // YAML 计算属性（标准 Kubernetes Pod YAML）
    const yamlContent = computed(() => {
      const podYaml = {
        apiVersion: "v1",
        kind: "Pod",
        metadata: {
          name: podDataCopy.name || "",
          namespace: podDataCopy.namespace || "default",
          labels: podDataCopy.labels || {}
        },
        spec: {
          containers: (podDataCopy.containers || []).map(c => ({
            name: c.name,
            image: c.image,
            imagePullPolicy: c.imagePullPolicy,
            ports: (c.ports || []).map(p => ({ containerPort: p }))
          }))
        }
      };
      try {
        return yaml.dump(podYaml);
      } catch {
        return "无法生成 YAML";
      }
    });

    // 复制 YAML
    const copyYaml = async () => {
      try {
        await navigator.clipboard.writeText(yamlContent.value);
        alert("已复制到剪贴板");
      } catch {
        alert("复制失败，请手动选择复制");
      }
    };

    return { podLabelsArray, podDataCopy, addLabel, removeLabel, syncLabels, updatePod, goBack, yamlContent, copyYaml };
  }
});
</script>

<style scoped>
.title > .icon {
  color: #4caf50;
}
.page-container {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  max-width: 1200px;
  margin: 0 auto;
}
.readonly-info {
  background: #f5f5f5;
  padding: 8px 12px;
  border-radius: 6px;
  margin-bottom: 15px;
  font-size: 14px;
}
.main-content {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
}
.form-content {
  flex: 2;
  min-width: 350px;
}
.yaml-preview {
  flex: 1;
  min-width: 300px;
  background: #fff;
  color: #333;
  padding: 16px;
  border-radius: 8px;
  font-size: 15px;
  position: relative;
  max-height: 600px;
  overflow-y: auto;
  border: 1px solid #eee;
}
.yaml-preview pre {
  white-space: pre-wrap;
  word-break: break-word;
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
.copy-btn {
  position: absolute;
  top: 8px;
  right: 12px;
  background: #42a5f5;
  color: #fff;
  border: none;
  padding: 4px 10px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
}
.form-section {
  margin-bottom: 20px;
  padding: 15px;
  border: 1px solid #eee;
  border-radius: 10px;
}
.section-icon.blue { color: #42a5f5; }
.section-icon.orange { color: #ff9800; }
.section-icon.purple { color: #9c27b0; }
.form-item {
  display: flex;
  align-items: center;
  margin: 10px 0;
  gap: 8px;
}
.form-item label {
  width: 80px;
  font-weight: 500;
  color: #555;
}
.form-item input,
.form-item textarea {
  flex: 1;
  padding: 6px 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
}
.container-card {
  background: #fafafa;
  border: 1px solid #ddd;
  border-radius: 10px;
  padding: 12px;
  margin-bottom: 12px;
}
.form-actions {
  margin-top: 20px;
  display: flex;
  gap: 12px;
}
.btn {
  padding: 6px 12px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
}
.btn-primary { background: #4caf50; color: #fff; }
.btn-secondary { background: #ccc; color: #333; }
.btn-add { background: #3178c6; color: #fff; margin-top: 10px; }
.btn-danger { background: #f66d6d; color: #fff; }
.btn-small { font-size: 12px; padding: 4px 8px; }
.labels-hint { font-size: 12px; color: #888; margin-bottom: 6px; }
</style>
