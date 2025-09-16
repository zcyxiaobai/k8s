<template>
  <div class="deployment-update-page">
    <div class="page-container">
      <h2 class="title">
        <i class="fa fa-edit icon"></i> 修改 Deployment - {{ depDataCopy?.name }}
      </h2>

      <!-- 只读基本信息 -->
      <div class="readonly-info">
        <p><strong>apiVersion：</strong> {{ depDataCopy?.apiVersion }}</p>
        <p><strong>kind：</strong> {{ depDataCopy?.kind }}</p>
        <p><strong>命名空间：</strong> {{ depDataCopy?.namespace }}</p>
      </div>

      <div class="main-content">
        <!-- 左侧表单 -->
        <div class="form-content">
          <!-- Labels 可编辑 -->
          <div class="form-section">
            <h3>
              <i class="fa fa-tag section-icon orange"></i> Deployment Labels
            </h3>
            <p class="labels-hint">⚠️ 可编辑 metadata.labels</p>
            <div
              v-for="(label, idx) in depLabelsArray"
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

          <!-- Replicas 可编辑 -->
          <div class="form-section">
            <h3>
              <i class="fa fa-copy section-icon blue"></i> 副本数
            </h3>
            <div class="form-item">
              <label>Replicas：</label>
              <input
                type="number"
                v-model.number="replicas"
                min="1"
                @input="syncReplicas"
              />
            </div>
          </div>

          <!-- Selector 只读 -->
          <div class="form-section">
            <h3>
              <i class="fa fa-crosshairs section-icon purple"></i> Selector (只读)
            </h3>
            <p class="labels-hint">⚠️ Deployment 的 selector 不可修改</p>
            <div class="selector-list">
              <div
                v-for="(value, key) in depDataCopy?.selector?.matchLabels || depDataCopy?.selector || {}"
                :key="'selector-' + key"
                class="selector-item"
              >
                <strong>{{ key }}:</strong> {{ value }}
              </div>
              <div v-if="!depDataCopy?.selector || Object.keys(depDataCopy?.selector || {}).length === 0">
                -
              </div>
            </div>
          </div>

          <!-- 容器信息 只读 -->
          <div class="form-section">
            <h3>
              <i class="fa fa-box section-icon purple"></i> 容器信息 (只读)
            </h3>
            <p class="labels-hint">⚠️ 容器规格不可修改，如需修改请新建 Deployment</p>
            <div
              v-for="(container, idx) in depDataCopy?.containers || []"
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
                  :value="formatPorts(container.ports)"
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
            <p class="labels-hint">⚠️ 只显示，无法修改</p>
            <div v-for="(value, key) in depDataCopy" :key="key" class="form-item">
              <template v-if="key !== 'labels' && key !== 'replicas' && key !== 'selector' && key !== 'containers'">
                <label>{{ key }}：</label>
                <input type="text" :value="value" readonly />
              </template>
            </div>
          </div>
        </div>

        <!-- 右侧 YAML 预览 -->
        <div class="yaml-preview">
          <h3><i class="fa fa-file-code-o"></i> Deployment 数据预览 (YAML)</h3>
          <button class="copy-btn" @click="copyYaml">复制</button>
          <pre>{{ yamlText }}</pre>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="form-actions">
        <button class="btn btn-primary" @click="updateDeployment">
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
import YAML from "yaml";

export default defineComponent({
  name: "UpdateDeployment",
  props: {
    depData: { type: Object, default: null }
  },
  setup(props) {
    const router = useRouter();
    const replicas = ref(props.depData?.replicas || 1);
    const depLabelsArray = ref([{ key: "", value: "" }]);

    // 复制 depData 用于实时更新
    const depDataCopy = reactive(JSON.parse(JSON.stringify(props.depData || {})));

    // 初始化 labels
    onMounted(() => {
      if (props.depData?.labels) {
        depLabelsArray.value = Object.entries(props.depData.labels).map(
          ([k, v]) => ({ key: k, value: v })
        );
      }
    });

    // 同步 Labels
    const syncLabels = () => {
      const labelsObj = {};
      depLabelsArray.value.forEach(l => {
        if (l.key) labelsObj[l.key] = l.value;
      });
      depDataCopy.labels = labelsObj;
    };

    // 同步 Replicas
    const syncReplicas = () => {
      depDataCopy.replicas = replicas.value;
    };

    const addLabel = () => depLabelsArray.value.push({ key: "", value: "" });
    const removeLabel = (idx) => {
      depLabelsArray.value.splice(idx, 1);
      syncLabels();
    };

    const updateDeployment = async () => {
  const labelsObj = {};
  depLabelsArray.value.forEach((l) => {
    if (l.key) labelsObj[l.key] = l.value;
  });

  try {
    // 调用新的接口 UpdateDev
    await axios.post("http://192.168.216.50:8090/dep/UpdateDev", {
      namespace: depDataCopy.namespace,
      name: depDataCopy.name,
      replicas: replicas.value,
      labels: labelsObj
    });
    alert("修改成功");
    router.back();
  } catch (err) {
    console.error(err);
    alert("修改失败，请检查控制台信息");
  }
};


    const goBack = () => router.back();

    const copyYaml = async () => {
      try {
        await navigator.clipboard.writeText(YAML.stringify(depDataCopy));
        alert("已复制 YAML 到剪贴板");
      } catch {
        alert("复制失败，请手动选择复制");
      }
    };

    // 格式化容器端口
    const formatPorts = (ports) => {
      if (!ports) return "-";
      if (Array.isArray(ports)) {
        if (ports.length === 0) return "-";
        if (typeof ports[0] === "number") return ports.join(",");
        if (typeof ports[0] === "object" && ports[0].containerPort) {
          return ports.map(p => p.containerPort).join(",");
        }
      }
      return "-";
    };

    // YAML 格式文本
    const yamlText = computed(() => YAML.stringify(depDataCopy));

    return {
      replicas,
      depLabelsArray,
      depDataCopy,
      addLabel,
      removeLabel,
      syncLabels,
      syncReplicas,
      updateDeployment,
      goBack,
      formatPorts,
      copyYaml,
      yamlText
    };
  }
});
</script>

<style scoped>
.page-container {
  max-width: 1000px;
  margin: 0 auto;
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 4px 10px rgba(0,0,0,0.1);
}

.title > .icon {
  color: #4caf50;
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
  gap: 20px;
}
.form-content {
  flex: 1;
}
.yaml-preview {
  flex: 1;
  background: #fff;
  color: #333;
  padding: 16px;
  border-radius: 8px;
  font-size: 15px;
  white-space: pre-wrap;
  word-wrap: break-word;
  overflow-x: hidden;
  position: relative;
  border: 2px solid transparent;
  background-clip: padding-box;
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
  max-width: 400px;
  padding: 6px 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-family: monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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
.btn-primary { background: #61ddaa; color: #fff; }
.btn-secondary { background: #ccc; color: #333; }
.btn-add { background: #3178c6; color: #fff; margin-top: 10px; }
.btn-danger { background: #f66d6d; color: #fff; }
.btn-small { font-size: 12px; padding: 4px 8px; }
.labels-hint { font-size: 12px; color: #888; margin-bottom: 6px; }

.selector-list { display: flex; flex-direction: column; gap: 4px; font-family: monospace; font-size: 13px; }
.selector-item { background: #f7f7f7; padding: 4px 8px; border-radius: 6px; word-break: break-all; }

.json-pre { white-space: pre-wrap; word-wrap: break-word; }
</style>
