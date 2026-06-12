    const { createApp, ref, computed, onMounted } = Vue;
    const mountId = 'app-[[.ID]]';
    const endpoint = '[[.Endpoint]]';

    const FlowNode = {
        name: 'flow-node',
        props: ['block', 'path', 'selectedPath', 'definitions'],
        template: `
            <div class="flow-block-container">
                <div class="block-card"
                     :class="{ active: isSelected }"
                     @click.stop="$emit('select', path)">

                    <div class="block-actions">
                        <button @click.stop="$emit('remove', path)" class="btn btn-danger btn-xs py-0 px-1" style="font-size: 0.6rem;">
                            <i class="bi bi-trash"></i>
                        </button>
                    </div>

                    <div class="d-flex align-items-center">
                        <div class="block-icon">
                            <i :class="getIcon(block.type)"></i>
                        </div>
                        <div class="overflow-hidden">
                            <div class="fw-bold small text-truncate">{{ block.title }}</div>
                            <div class="text-muted" style="font-size: 0.65rem;">{{ block.type }}</div>
                        </div>
                    </div>
                </div>

                <!-- Branches -->
                        <template v-if="hasBranches">
                            <div class="branches-container">
                                <template v-for="branchName in branchNames" :key="branchName">
                                    <div class="branch">
                                        <div class="connector-v"></div>
                                        <div class="branch-label">{{ branchName }}</div>
                                        <div class="connector-v"></div>

                                        <template v-for="(b, idx) in (block.branches[branchName] || [])" :key="b.id">
                                            <flow-node
                                                :block="b"
                                                :path="[...path, 'branches', branchName, idx]"
                                                :selected-path="selectedPath"
                                                :definitions="definitions"
                                                @select="(p) => $emit('select', p)"
                                                @remove="(p) => $emit('remove', p)"
                                                @add-to-branch="(p, def) => $emit('add-to-branch', p, def)"
                                            ></flow-node>
                                            <div v-if="idx < block.branches[branchName].length - 1" class="connector-v"></div>
                                        </template>

                                        <button @click.stop="addBlockToThisBranch(branchName)" class="btn btn-outline-primary btn-xs mt-2" style="font-size: 0.6rem; padding: 2px 5px;">
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
                const def = this.definitions.find(d => d.type === this.block.type);
                return def && def.branchNames && def.branchNames.length > 0;
            },
            branchNames() {
                const def = this.definitions.find(d => d.type === this.block.type);
                return def ? def.branchNames : [];
            }
        },
        methods: {
            getIcon(type) {
                const def = this.definitions.find(d => d.type === type);
                return def ? def.icon : 'bi-box';
            },
            addBlockToThisBranch(branchName) {
                this.$emit('add-to-branch', [...this.path, 'branches', branchName], null);
            }
        }
    };

    createApp({
        components: { FlowNode },
        setup() {
            const flow = ref([]);
            const definitions = ref([]);
            const selectedBlockPath = ref(null);
            const saving = ref(false);

            const selectedBlock = computed(() => {
                if (!selectedBlockPath.value) return null;
                let curr = flow.value;
                for (const p of selectedBlockPath.value) {
                    curr = curr[p];
                }
                return curr;
            });

            const addBlock = (def) => {
                const newBlock = {
                    id: 'b_' + Math.random().toString(36).substr(2, 9),
                    type: def.type,
                    title: def.title,
                    data: JSON.parse(JSON.stringify(def.defaultData || {})),
                    branches: {}
                };
                if (def.branchNames) {
                    def.branchNames.forEach(bn => newBlock.branches[bn] = []);
                }
                flow.value.push(newBlock);
            };

            const addBlockToBranch = (branchPath, def) => {
                if (!def) def = definitions.value[0];

                const newBlock = {
                    id: 'b_' + Math.random().toString(36).substr(2, 9),
                    type: def.type,
                    title: def.title,
                    data: JSON.parse(JSON.stringify(def.defaultData || {})),
                    branches: {}
                };
                if (def.branchNames) {
                    def.branchNames.forEach(bn => newBlock.branches[bn] = []);
                }

                let curr = flow.value;
                for (const p of branchPath) {
                    curr = curr[p];
                }
                curr.push(newBlock);
            };

            const selectBlock = (path) => {
                selectedBlockPath.value = path;
            };

            const removeBlock = (path) => {
                const idx = path[path.length - 1];
                const parentPath = path.slice(0, -1);
                let curr = flow.value;
                for (const p of parentPath) {
                    curr = curr[p];
                }
                curr.splice(idx, 1);
                selectedBlockPath.value = null;
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
                    if (globalData.definitions) definitions.value = globalData.definitions;
                }

                // Then refresh from API to ensure we have latest (optional)
                try {
                    const baseUrl = endpoint.endsWith('/') ? endpoint : endpoint + '/';
                    const res = await fetch(baseUrl + 'config');
                    const data = await res.json();
                    definitions.value = data.definitions;
                    flow.value = data.flow || [];
                } catch (e) {
                    console.warn("Failed to refresh config from API", e);
                }
            };

            onMounted(loadConfig);

            return {
                flow,
                definitions,
                selectedBlockPath,
                selectedBlock,
                saving,
                addBlock,
                addBlockToBranch,
                selectBlock,
                removeBlock,
                saveFlow
            };
        }
    }).mount('#' + mountId);
