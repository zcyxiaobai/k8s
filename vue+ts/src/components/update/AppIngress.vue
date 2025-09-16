<template>
  <div class="ingress-update-page" v-if="ingData">
    <div class="page-container">
      <h2 class="title">
        <i class="fa fa-edit icon"></i> 修改 Ingress - {{ ingress.name }}
      </h2>

      <!-- 只读固定字段 -->
      <div class="readonly-info">
        <p><strong>apiVersion：</strong> networking.k8s.io/v1</p>
        <p><strong>kind：</strong> Ingress</p>
        <p><strong>命名空间：</strong> {{ ingress.namespace }}</p>
        <p v-if="ingress.ingressClassName"><strong>IngressClass：</strong> {{ ingress.ingressClassName }}</p>
      </div>

      <div class="main-content">
        <!-- 左边表单 -->
        <div class="form-content">
          <!-- Labels -->
          <div class="form-section">
            <h3><i class="fa fa-tag section-icon purple"></i> Labels</h3>
            <div
              v-for="(label, idx) in ingress.labels"
              :key="'label-' + idx"
              class="form-item"
            >
              <label>Key：</label>
              <input
                v-model="label.key"
                type="text"
                :readonly="label.readonly"
                :title="label.key"
                placeholder="如 app"
              />
              <label>Value：</label>
              <input
                v-model="label.value"
                type="text"
                :readonly="label.readonly"
                :title="label.value"
                placeholder="如 nginx"
              />
              <button v-if="!label.readonly" class="btn btn-danger btn-small" @click="removeLabel(idx)">
                <i class="fa fa-trash"></i>
              </button>
            </div>
            <button class="btn btn-add" @click="addLabel">
              <i class="fa fa-plus"></i> 添加 Label
            </button>
          </div>

          <!-- Annotations -->
          <div class="form-section">
            <h3><i class="fa fa-sticky-note section-icon orange"></i> Annotations</h3>
            <div
              v-for="(annotation, idx) in ingress.annotations"
              :key="'anno-' + idx"
              class="form-item"
            >
              <label>Key：</label>
              <input
                v-model="annotation.key"
                type="text"
                :readonly="annotation.readonly"
                :title="annotation.key"
                placeholder="如 nginx.ingress.kubernetes.io/rewrite-target"
              />
              <label>Value：</label>
              <input
                v-model="annotation.value"
                type="text"
                :readonly="annotation.readonly"
                :title="annotation.value"
                placeholder="如 /"
              />
              <button v-if="!annotation.readonly" class="btn btn-danger btn-small" @click="removeAnnotation(idx)">
                <i class="fa fa-trash"></i>
              </button>
            </div>
            <button class="btn btn-add" @click="addAnnotation">
              <i class="fa fa-plus"></i> 添加 Annotation
            </button>
          </div>

          <!-- Rules 配置 -->
          <div class="form-section">
            <h3><i class="fa fa-link section-icon green"></i> Rules 配置</h3>
            <div
              v-for="(rule, idx) in ingress.rules"
              :key="'rule-' + idx"
              class="rule-item"
            >
              <div class="form-item">
                <label>Host：</label>
                <input
                  v-model="rule.host"
                  type="text"
                  :readonly="rule.readonly"
                  :title="rule.host"
                  placeholder="如 nginx.example.com"
                />
              </div>

              <div
                v-for="(pathObj, pIdx) in rule.paths"
                :key="'path-' + pIdx"
                class="path-item"
              >
                <div class="path-row">
                  <label>路径：</label>
                  <input
                    v-model="pathObj.path"
                    type="text"
                    :readonly="pathObj.readonly"
                    :title="pathObj.path"
                    placeholder="如 /"
                  />
                </div>
                <div class="path-row">
                  <label>PathType：</label>
                  <select v-model="pathObj.pathType" :disabled="pathObj.readonly">
                    <option value="Prefix">Prefix</option>
                    <option value="Exact">Exact</option>
                    <option value="ImplementationSpecific">ImplementationSpecific</option>
                  </select>
                </div>
                <div class="path-row">
                  <label>服务：</label>
                  <input
                    v-model="pathObj.serviceName"
                    type="text"
                    :readonly="pathObj.readonly"
                    :title="pathObj.serviceName"
                    placeholder="如 nginx-service"
                  />
                </div>
                <div class="path-row">
                  <label>端口：</label>
                  <input
                    v-model="pathObj.servicePort"
                    type="number"
                    :readonly="pathObj.readonly"
                    :title="pathObj.servicePort"
                    placeholder="如 80"
                  />
                </div>
                <button v-if="!pathObj.readonly" class="btn btn-danger btn-small" @click="removePath(idx, pIdx)">
                  <i class="fa fa-trash"></i> 移除路径
                </button>
              </div>
              <button v-if="!rule.readonly" class="btn btn-add btn-small" @click="addPath(idx)">
                <i class="fa fa-plus"></i> 添加路径
              </button>
              <button v-if="!rule.readonly" class="btn btn-danger btn-small" style="margin-top:6px" @click="removeRule(idx)">
                <i class="fa fa-trash"></i> 移除 Host
              </button>
            </div>
            <button class="btn btn-add" @click="addRule">
              <i class="fa fa-plus"></i> 添加 Host
            </button>
          </div>

          <!-- 按钮 -->
          <div class="form-actions">
            <button class="btn btn-primary" @click="updateIngress">
              <i class="fa fa-check"></i> 保存修改
            </button>
            <button class="btn btn-secondary" @click="goBack">
              <i class="fa fa-refresh"></i> 取消
            </button>
          </div>
        </div>

        <!-- 右边 YAML 实时预览 -->
        <div class="yaml-preview">
          <h3><i class="fa fa-file-code-o"></i> YAML 预览</h3>
          <button class="copy-btn" @click="copyYaml">复制</button>
          <pre class="yaml-content">{{ yamlOutput }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
// import axios from "axios";
import axios from "@/main" 
import yaml from "js-yaml";

export default defineComponent({
  name: "AppIngressUpdate",
  props: { ingData: { type: Object, default: null } },
  setup(props) {
    const router = useRouter();
    const ingress = ref({ name: "", namespace: "", ingressClassName: "", labels: [], annotations: [], rules: [] });

    onMounted(() => {
      if (props.ingData) {
        const data = props.ingData;
        ingress.value.name = data.name || "";
        ingress.value.namespace = data.namespace || "";
        ingress.value.ingressClassName = data.ingressClassName || "";
        ingress.value.labels = Object.entries(data.labels || {}).map(([k, v]) => ({ key: k, value: v, readonly: false }));
        ingress.value.annotations = Object.entries(data.annotations || {}).map(([k, v]) => ({ key: k, value: v, readonly: false }));
        ingress.value.rules = (data.rules || []).map(r => ({
          host: r.host || "",
          readonly: false,
          paths: (r.paths || []).map(p => ({
            path: p.path || "",
            pathType: p.pathType || "Prefix",
            serviceName: p.service || "",
            servicePort: p.port || 80,
            readonly: false
          }))
        }));
        if (ingress.value.labels.length === 0) ingress.value.labels.push({ key: "", value: "", readonly: false });
        if (ingress.value.annotations.length === 0) ingress.value.annotations.push({ key: "", value: "", readonly: false });
        if (ingress.value.rules.length === 0) ingress.value.rules.push({ host: "", readonly: false, paths: [{ path: "", pathType: "Prefix", serviceName: "", servicePort: 80, readonly: false }] });
      }
    });

    const addLabel = () => ingress.value.labels.push({ key: "", value: "", readonly: false });
    const removeLabel = i => ingress.value.labels.splice(i, 1);
    const addAnnotation = () => ingress.value.annotations.push({ key: "", value: "", readonly: false });
    const removeAnnotation = i => ingress.value.annotations.splice(i, 1);
    const addRule = () => ingress.value.rules.push({ host: "", readonly: false, paths: [{ path: "", pathType: "Prefix", serviceName: "", servicePort: 80, readonly: false }] });
    const removeRule = i => ingress.value.rules.splice(i, 1);
    const addPath = ruleIdx => ingress.value.rules[ruleIdx].paths.push({ path: "", pathType: "Prefix", serviceName: "", servicePort: 80, readonly: false });
    const removePath = (ruleIdx, pathIdx) => ingress.value.rules[ruleIdx].paths.splice(pathIdx, 1);

    const yamlOutput = computed(() => {
      const obj = {
        apiVersion: "networking.k8s.io/v1",
        kind: "Ingress",
        metadata: {
          name: ingress.value.name,
          namespace: ingress.value.namespace,
          labels: ingress.value.labels.reduce((acc, cur) => { if(cur.key) acc[cur.key]=cur.value; return acc; }, {}),
          annotations: ingress.value.annotations.reduce((acc, cur) => { if(cur.key) acc[cur.key]=cur.value; return acc; }, {})
        },
        spec: {
          ingressClassName: ingress.value.ingressClassName || undefined,
          rules: ingress.value.rules.map(r => ({ host: r.host, http: { paths: r.paths.map(p => ({ path: p.path, pathType: p.pathType, backend: { service: { name: p.serviceName, port: { number: Number(p.servicePort) || 80 } } } })) } }))
        }
      };
      return yaml.dump(obj, { noRefs: true });
    });

    const updateIngress = async () => {
      const payload = {
        name: ingress.value.name,
        namespace: ingress.value.namespace,
        labels: ingress.value.labels.reduce((acc, cur) => { if(cur.key) acc[cur.key]=cur.value; return acc; }, {}),
        annotations: ingress.value.annotations.reduce((acc, cur) => { if(cur.key) acc[cur.key]=cur.value; return acc; }, {}),
        rules: ingress.value.rules.map(r => ({ host: r.host, paths: r.paths.map(p => ({ path: p.path, pathType: p.pathType, serviceName: p.serviceName, servicePort: p.servicePort })) }))
      };
      try { const res = await axios.post("http://192.168.216.50:8090/ingest/UpdateIngress", payload); alert(res.data.message || "修改成功"); router.back(); }
      catch(err) { console.error(err); alert("修改失败，请检查控制台"); }
    };

    const copyYaml = async () => { try { await navigator.clipboard.writeText(yamlOutput.value); alert("已复制到剪贴板"); } catch { alert("复制失败，请手动复制"); } };
    const goBack = () => router.back();

    return { ingress, addLabel, removeLabel, addAnnotation, removeAnnotation, addRule, removeRule, addPath, removePath, updateIngress, yamlOutput, copyYaml, goBack };
  }
});
</script>

<style scoped>
.title>.icon{color:#9c27b0;}
.page-container{background:#fff;border-radius:12px;padding:20px;box-shadow:0 4px 10px rgba(0,0,0,.1);}
.readonly-info{background:#f5f5f5;padding:8px 12px;border-radius:6px;margin-bottom:15px;font-size:14px;}
.main-content{display:flex;gap:20px;}
.form-content{flex:1.5;} /* 输入框增加占比，拉宽 */
.yaml-preview{
  flex:1;
  background:#fff;
  color:#333;
  padding:16px;
  border-radius:8px;
  font-size:15px;
  overflow-x:auto;
  position:relative;
  border:2px solid transparent;
  background-clip: padding-box;
}
.yaml-preview::before {
  content:"";
  position:absolute;
  inset:0;
  border-radius:8px;
  padding:2px;
  background: linear-gradient(45deg, #42a5f5, #66bb6a);
  -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
  mask-composite: exclude;
}
.yaml-content{word-break:break-word;white-space:pre-wrap;overflow-wrap:break-word;}
.copy-btn{position:absolute;top:8px;right:12px;background:#42a5f5;color:#fff;border:none;padding:4px 10px;border-radius:6px;cursor:pointer;font-size:12px;}
.form-section{margin-bottom:20px;padding:15px;border:1px solid #eee;border-radius:10px;}
.section-icon.purple{color:#9c27b0;}
.section-icon.green{color:#66bb6a;}
.section-icon.orange{color:#ff9800;}
.form-item{display:flex;align-items:center;margin:10px 0;gap:8px;flex-wrap:wrap;}
.form-item label{width:150px;font-weight:500;color:#555;flex-shrink:0;} /* 标签宽度略减 */
.form-item input,.form-item select{flex:1.2;min-width:120px;padding:6px 10px;border:1px solid #ddd;border-radius:6px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;} /* 输入框自适应宽度 */
.rule-item{flex-direction:column;border:1px dashed #ddd;padding:10px;border-radius:8px;margin-bottom:10px;}
.path-item{display:flex;flex-direction:column;gap:6px;margin-bottom:6px;}
.path-row{display:flex;flex-wrap:wrap;gap:8px;align-items:center;}
.form-actions{margin-top:20px;display:flex;gap:12px;}
.btn{padding:6px 12px;border:none;border-radius:8px;cursor:pointer;font-size:14px;}
.btn-primary{background:#9c27b0;color:#fff;}
.btn-secondary{background:#ccc;color:#333;}
.btn-add{background:#3178c6;color:#fff;margin-top:10px;}
.btn-danger{background:#f66d6d;color:#fff;}
.btn-small{font-size:12px;padding:4px 8px;}
</style>
