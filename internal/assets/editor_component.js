    const { createApp, ref, computed, onMounted } = Vue;
    const mountId = 'app-[[.ID]]';
    const endpoint = '[[.Endpoint]]';

    const StepNode = {
        name: 'step-node',
        props: ['step', 'path', 'selectedPath', 'stepDefinitions'],
        template: `
            <div class="FlowStepContainer">
                <div class="StepCard"
                     :class="{ active: isSelected }"
                     @click.stop="$emit('select', path)">

                    <div class="StepActions">
                        <button @click.stop="$emit('remove', path)" class="btn btn-danger btn-xs py-0 px-1" style="font-size: 0.6rem;">
                            <i class="bi bi-trash"></i>
                        </button>
                    </div>

                    <div class="d-flex align-items-center">
                        <div class="StepIcon">
                            <i :class="getIcon(step.type)"></i>
                        </div>
                        <div class="overflow-hidden">
                            <div class="fw-bold small text-truncate">{{ step.title }}</div>
                            <div class="text-muted" style="font-size: 0.65rem;">{{ step.type }}</div>
                        </div>
                    </div>
                </div>

                <!-- Branches -->
                        <template v-if="hasBranches">
                            <div class="BranchesContainer">
                                <template v-for="branchName in branchNames" :key="branchName">
                                    <div class="Branch">
                                        <div class="ConnectorV"></div>
                                        <div class="BranchLabel">{{ branchName }}</div>
                                        <div class="ConnectorV"></div>

                                        <template v-for="(s, idx) in (step.branches[branchName] || [])" :key="s.id">
                                            <step-node
                                                :step="s"
                                                :path="[...path, 'branches', branchName, idx]"
                                                :selected-path="selectedPath"
                                                :step-definitions="stepDefinitions"
                                                @select="(p) => $emit('select', p)"
                                                @remove="(p) => $emit('remove', p)"
                                                @add-to-branch="(p, def) => $emit('add-to-branch', p, def)"
                                            ></step-node>
                                            <div v-if="idx < step.branches[branchName].length - 1" class="ConnectorV"></div>
                                        </template>

                                        <button @click.stop="addStepToThisBranch(branchName)" class="btn btn-outline-primary btn-xs mt-2" style="font-size: 0.6rem; padding: 2px 5px;">
                                            <i class="bi bi-plus"></i>
                                        </button>
                                    </div>
                                </template>
                    </div>
                        </template>
            </div>
        `,
        computed: {
            isSelected() {
                return JSON.stringify(this.path) === JSON.stringify(this.selectedPath);
            },
            hasBranches() {
                const def = this.stepDefinitions.find(d => d.type === this.step.type);
                return def && def.branchNames && def.branchNames.length > 0;
            },
            branchNames() {
                const def = this.stepDefinitions.find(d => d.type === this.step.type);
                return def ? def.branchNames : [];
            }
        },
        methods: {
            getIcon(type) {
                const def = this.stepDefinitions.find(d => d.type === type);
                return def ? def.icon : 'bi-box';
            },
            addStepToThisBranch(branchName) {
                this.$emit('add-to-branch', [...this.path, 'branches', branchName], null);
            }
        }
    };

    createApp({
        components: { StepNode },
        setup() {
            const flow = ref([]);
            const stepDefinitions = ref([]);
            const selectedStepPath = ref(null);
            const saving = ref(false);

            const selectedStep = computed(() => {
                if (!selectedStepPath.value) return null;
                let curr = flow.value;
                for (const p of selectedStepPath.value) {
                    curr = curr[p];
                }
                return curr;
            });

            const addStep = (def) => {
                const newStep = {
                    id: 's_' + Math.random().toString(36).substr(2, 9),
                    type: def.type,
                    title: def.title,
                    data: JSON.parse(JSON.stringify(def.defaultData || {})),
                    branches: {}
                };
                if (def.branchNames) {
                    def.branchNames.forEach(bn => newStep.branches[bn] = []);
                }
                flow.value.push(newStep);
            };

            const addStepToBranch = (branchPath, def) => {
                if (!def) def = stepDefinitions.value[0];

                const newStep = {
                    id: 's_' + Math.random().toString(36).substr(2, 9),
                    type: def.type,
                    title: def.title,
                    data: JSON.parse(JSON.stringify(def.defaultData || {})),
                    branches: {}
                };
                if (def.branchNames) {
                    def.branchNames.forEach(bn => newStep.branches[bn] = []);
                }

                let curr = flow.value;
                for (const p of branchPath) {
                    curr = curr[p];
                }
                curr.push(newStep);
            };

            const selectStep = (path) => {
                selectedStepPath.value = path;
            };

            const removeStep = (path) => {
                const idx = path[path.length - 1];
                const parentPath = path.slice(0, -1);
                let curr = flow.value;
                for (const p of parentPath) {
                    curr = curr[p];
                }
                curr.splice(idx, 1);
                selectedStepPath.value = null;
            };

            const saveFlow = async () => {
                saving.value = true;
                try {
                    const baseUrl = endpoint.endsWith('/') ? endpoint : endpoint + '/';
                    await fetch(baseUrl + 'save', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(flow.value)
                    });
                    alert("Flow published successfully!");
                } catch (e) {
                    console.error(e);
                    alert("Save failed");
                } finally {
                    saving.value = false;
                }
            };

            const loadConfig = async () => {
                // Pre-load from global if available (injected via HTML template)
                const globalData = window["editorData_[[.ID]]"];
                if (globalData) {
                    if (globalData.flow) flow.value = globalData.flow;
                    if (globalData.stepDefinitions) stepDefinitions.value = globalData.stepDefinitions;
                }

                // Then refresh from API to ensure we have latest (optional)
                try {
                    const baseUrl = endpoint.endsWith('/') ? endpoint : endpoint + '/';
                    const res = await fetch(baseUrl + 'config');
                    const data = await res.json();
                    stepDefinitions.value = data.stepDefinitions;
                    flow.value = data.flow || [];
                } catch (e) {
                    console.warn("Failed to refresh config from API", e);
                }
            };

            onMounted(loadConfig);

            return {
                flow,
                stepDefinitions,
                selectedStepPath,
                selectedStep,
                saving,
                addStep,
                addStepToBranch,
                selectStep,
                removeStep,
                saveFlow
            };
        }
    }).mount('#' + mountId);
