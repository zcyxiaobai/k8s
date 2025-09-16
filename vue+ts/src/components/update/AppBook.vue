<template>
  <div class="bookapp-update-page">
    <div class="page-container">
      <h2 class="title">
        <i class="fa fa-edit icon"></i> 修改 Bookapp
      </h2>

      <!-- 只读基本信息 -->
      <div class="readonly-info">
        <p><strong>apiVersion：</strong> webapp.y2505.com/v1</p>
        <p><strong>kind：</strong> Bookapp</p>
        <p><strong>名称：</strong> {{ bookapp.name }}</p>
        <p><strong>命名空间：</strong> {{ bookapp.namespace }}</p>
      </div>

      <div class="main-content">
        <!-- 左边表单 -->
        <div class="form-content">
          <!-- Labels -->
          <div class="form-section">
            <h3>
              <i class="fa fa-tag section-icon orange"></i> Labels
            </h3>
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
              <input type="number" v-model="bookapp.spec.size" />
            </div>
            <div class="form-item">
              <label>镜像：</label>
              <input type="text" v-model="bookapp.spec.image" />
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
              <input type="text" v-model="bookapp.spec.port.name" />
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
              <input type="number" v-model="bookapp.spec.port.port" />
            </div>
            <div class="form-item">
              <label>目标端口：</label>
              <input type="number" v-model="bookapp.spec.port.targetPort" />
            </div>
            <div class="form-item">
              <label>域名：</label>
              <input type="text" v-model="bookapp.spec.serverName" />
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
        <button class="btn btn-primary" @click="updateBookapp">
          <i class="fa fa-check"></i> 保存修改
        </button>
        <button class="btn btn-secondary" @click="cancelForm">
          <i class="fa fa-refresh"></i> 取消
        </button>
       
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref, computed, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import yaml from "js-yaml";
// import axios from "axios";
import axios from "@/main" 

export default defineComponent({
  name: "AppBookappUpdate",
  setup() {
    const route = useRoute();
    const router = useRouter();

    const bookapp = ref({
      name: "",
      namespace: "",
      labels: {},
      spec: { size: 1, image: "", imagePullPolicy: "IfNotPresent", port: {}, serverName: "" }
    //   spec: { size: 1, image: "", imagePullPolicy: "IfNotPresent", port: { name:"", protocol:"TCP", port: 80, targetPort: 80 }, serverName: "" }

    });

    const labelsArray = ref([]);

    // ✅ 初始化：从路由 query 取 bookapp 数据
    onMounted(() => {
      if (route.query.bookapp) {
        try {
          const parsed = JSON.parse(decodeURIComponent(route.query.bookapp));
          bookapp.value = parsed;

          // 还原 labels 为数组形式
          if (parsed.labels) {
            labelsArray.value = Object.entries(parsed.labels).map(([k, v]) => ({ key: k, value: v }));
          }
        } catch (e) {
          console.error("解析 bookapp 参数失败:", e);
        }
      }
    });
const cancelForm = () => {
  router.push("/welcome/crd-view")   // 跳回 BookApp 列表页
}
    const updateLabels = () => {
      const obj = {};
      labelsArray.value.forEach(l => { if (l.key) obj[l.key] = l.value });
      bookapp.value.labels = obj;
    };
    const addLabel = () => labelsArray.value.push({ key: "", value: "" });
    const removeLabel = idx => { labelsArray.value.splice(idx, 1); updateLabels() };

    const yamlOutput = computed(() => yaml.dump(bookapp.value, { noRefs: true }));

    const copyYaml = async () => {
      try { await navigator.clipboard.writeText(yamlOutput.value); alert("已复制到剪贴板"); }
      catch { alert("复制失败"); }
    };

    const resetForm = () => {
      // 恢复路由传过来的数据
      if (route.query.bookapp) {
        const parsed = JSON.parse(decodeURIComponent(route.query.bookapp));
        bookapp.value = parsed;
        labelsArray.value = Object.entries(parsed.labels || {}).map(([k, v]) => ({ key: k, value: v }));
      }
    };

    const updateBookapp = async () => {
      try {
        const payload = {
          namespace: bookapp.value.namespace,
          name: bookapp.value.name,
          labels: bookapp.value.labels,
          spec: bookapp.value.spec
        };

        // const res = await axios.post("http://192.168.216.50:8090/bookapp/updatebookapp", payload);
        const res = await axios.post("http://192.168.216.50:8090/bookapp/UpdateBookApp", payload);

        alert(res.data?.message || "✅ 修改成功");
       router.push("/welcome/crd-view")   // 跳回 BookApp 列表页
      } catch (err) {
        console.error("修改失败:", err);
        alert("❌ 修改失败，请检查后台日志");
      }
    };

    return {
      bookapp, labelsArray, yamlOutput,
      updateLabels, addLabel, removeLabel,
      copyYaml, resetForm, updateBookapp,cancelForm
    };
  }
});
</script>

<style scoped>
.page-container { background: #fff; border-radius: 12px; padding: 20px; box-shadow: 0 4px 10px rgba(0,0,0,0.1); }
.title>.icon { color: #42a5f5 }
.readonly-info { background: #f5f5f5; padding: 8px 12px; border-radius: 6px; margin-bottom: 15px; font-size: 14px }
.main-content { display: flex; gap: 20px }
.form-content { flex: 1 }
.yaml-preview { flex: 1; background: #fff; padding: 16px; border-radius: 8px; font-size: 15px; overflow-x: auto; position: relative; border: 2px solid transparent; background-clip: padding-box; }
.yaml-preview::before { content: ""; position: absolute; inset: 0; border-radius: 8px; padding: 2px; background: linear-gradient(45deg, #42a5f5, #66bb6a); -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0); -webkit-mask-composite: xor; mask-composite: exclude; }
.copy-btn { position: absolute; top: 8px; right: 12px; background: #42a5f5; color: #fff; border: none; padding: 4px 10px; border-radius: 6px; cursor: pointer; font-size: 12px }
.form-section { margin-bottom: 20px; padding: 15px; border: 1px solid #eee; border-radius: 10px }
.section-icon.orange { color: #ff9800 }
.section-icon.purple { color: #9c27b0 }
.form-item { display: flex; align-items: center; margin: 10px 0; gap: 8px }
.form-item label { width: 80px; font-weight: 500; color: #555 }
.form-item input,
.form-item select {
  flex: unset;              /* 不要强制拉伸 */
  width: 100%;              /* 宽度自适应父容器 */
  max-width: 300px;         /* ✅ 限制最大宽度 */
  padding: 6px 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  box-sizing: border-box;   /* 避免 padding 撑大 */
}
.yaml-preview pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  max-width: 100%;
  overflow-x: auto;
}
.form-actions { margin-top: 20px; display: flex; gap: 12px }
.btn { padding: 6px 12px; border: none; border-radius: 8px; cursor: pointer; font-size: 14px }
.btn-primary { background: #42a5f5; color: #fff }
.btn-secondary { background: #ccc; color: #333 }
.btn-add { background: #3178c6; color: #fff; margin-top: 10px }
.btn-danger { background: #f66d6d; color: #fff }
.btn-small { font-size: 12px; padding: 4px 8px }
</style>
