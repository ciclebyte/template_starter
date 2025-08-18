<template>
  <n-drawer v-model:show="visible" :width="'90%'" placement="bottom" :height="drawerHeight + 'px'">
    <n-drawer-content title="模板变量定义" :native-scrollbar="false">
      <!-- 头部操作按钮 -->
      <div class="header-actions" style="display: flex; justify-content: flex-end; margin-bottom: 16px; gap: 8px;">
        <n-button size="small" @click="handleSaveClick" type="primary" :loading="saving">
          <template #icon>
            <n-icon><SaveOutline /></n-icon>
          </template>
          保存定义
        </n-button>
        <n-button size="small" @click="showTestDataModal" type="info" quaternary>
          <template #icon>
            <n-icon><FlaskOutline /></n-icon>
          </template>
          测试数据
        </n-button>
        <n-button size="small" @click="showVariableAnalysisModal" type="success" quaternary>
          <template #icon>
            <n-icon><SearchOutline /></n-icon>
          </template>
          分析变量
        </n-button>
        <n-button size="small" @click="emergencyCleanup" type="warning" quaternary>
          <template #icon>
            <n-icon><TrashOutline /></n-icon>
          </template>
          重置所有
        </n-button>
        <n-button size="small" quaternary @click="visible = false">
          <template #icon>
            <n-icon><CloseOutline /></n-icon>
          </template>
          关闭
        </n-button>
      </div>

      <!-- 拖拽手柄 -->
      <div class="resize-handle" @mousedown="startResize">
        <div class="handle-bar"></div>
      </div>

      <div class="expose-layout">
        <!-- 左侧：变量树 -->
        <div class="left-panel" :style="{ width: leftPanelWidth + 'px' }">
          <div class="panel-title">
            变量资源
          </div>
          <div class="variable-tree" @contextmenu="onTreeAreaContextMenu">
            <n-tree
              :data="variableTreeData"
              :selected-keys="selectedKeys"
              :expanded-keys="expandedKeys"
              key-field="key"
              label-field="title"
              children-field="children"
              :node-props="nodeProps"
              :render-switcher-icon="renderSwitcherIcon"
              :render-label="renderLabel"
              block-line
              selectable
              @update:selected-keys="onSelectVariable"
              @update:expanded-keys="onExpandKeys"
            >
              
              <template #suffix="{ option }">
                <div class="node-actions" @click.stop>
                  <n-dropdown
                    :options="getNodeMenuOptions(option)"
                    @select="(key) => handleNodeAction(key, option)"
                    trigger="click"
                    placement="bottom-end"
                  >
                    <n-button size="tiny" quaternary circle>
                      <template #icon>
                        <n-icon><EllipsisHorizontalOutline /></n-icon>
                      </template>
                    </n-button>
                  </n-dropdown>
                </div>
              </template>
            </n-tree>
            
            <!-- 空状态提示 -->
            <div v-if="!variableTreeData || variableTreeData.length === 0" 
                 class="empty-tree-hint"
                 @contextmenu="onTreeAreaContextMenu">
              暂无变量（右键添加）
            </div>
            
            <!-- 右键上下文菜单 -->
            <n-dropdown
              to="body"
              trigger="manual"
              :x="contextMenuX"
              :y="contextMenuY"
              :options="contextMenuOptions"
              :show="showContextMenuFlag"
              @select="handleContextMenuAction"
              @clickoutside="hideContextMenu"
              placement="bottom-start"
            />
          </div>
        </div>

        <!-- 左右分割线 -->
        <div class="panel-resizer left-resizer" @mousedown="startLeftResize"></div>

        <!-- 中间：变量信息 -->
        <div class="center-panel" :style="{ width: centerPanelWidth + 'px' }">
          <div class="panel-title">变量信息</div>
          
          <!-- 变量编辑表单 -->
          <div v-if="selectedVariableData" class="variable-form" :key="selectedVariableData.path">
            <div class="form-header">
              <span class="var-path">{{ selectedVariableData.path }}</span>
              <n-tag :type="getTypeTagType(selectedVariableData.type)" size="small">
                {{ getTypeLabel(selectedVariableData.type) }}
              </n-tag>
            </div>
            
            <n-form ref="formRef" :model="selectedVariableData" size="small" label-placement="top" :key="selectedVariableData.path + '_form'">
              <n-grid :cols="2" :x-gap="16" :y-gap="12">
                <!-- 基础信息 -->
                <n-grid-item :span="2">
                  <div class="form-section-title">基础信息</div>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="显示名称 (title)">
                    <n-input v-model:value="selectedVariableData.title" placeholder="用户友好的显示名称" />
                  </n-form-item>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="数据类型 (type)">
                    <n-select 
                      v-model:value="selectedVariableData.type" 
                      :options="typeOptions"
                      @update:value="onTypeChange"
                    />
                  </n-form-item>
                </n-grid-item>
                
                <n-grid-item :span="2">
                  <n-form-item label="描述信息 (description)">
                    <n-input 
                      v-model:value="selectedVariableData.description" 
                      type="textarea" 
                      :rows="2"
                      placeholder="一句话说明变量的用途"
                    />
                  </n-form-item>
                </n-grid-item>
                
                <n-grid-item :span="2">
                  <n-form-item label="插入文本 (insertText)">
                    <n-input 
                      v-model:value="selectedVariableData.insertText" 
                      placeholder="{{.变量名}}"
                    />
                  </n-form-item>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="是否必填 (required)">
                    <n-switch v-model:value="selectedVariableData.required">
                      <template #checked>必填</template>
                      <template #unchecked>可选</template>
                    </n-switch>
                  </n-form-item>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="默认值 (default)">
                    <!-- 对象、数组和对象数组类型显示JSON字符串 -->
                    <n-input 
                      v-if="selectedVariableData.type === 'object' || selectedVariableData.type === 'array' || selectedVariableData.type === 'object_arr'"
                      :value="JSON.stringify(selectedVariableData.default)"
                      :placeholder="`默认${getTypeLabel(selectedVariableData.type)}值`"
                      readonly
                      disabled
                    />
                    <!-- 其他类型正常处理 -->
                    <component
                      v-else
                      :is="getInputComponent(selectedVariableData.type)"
                      v-model:value="selectedVariableData.default"
                      :placeholder="`默认${getTypeLabel(selectedVariableData.type)}值`"
                    />
                  </n-form-item>
                </n-grid-item>

                <!-- 特殊类型配置 -->
                <template v-if="selectedVariableData.type === 'enum'">
                  <n-grid-item :span="2">
                    <div class="form-section-title">枚举配置</div>
                  </n-grid-item>
                  <n-grid-item :span="2">
                    <n-form-item label="可选值 (enum)">
                      <n-dynamic-tags 
                        v-model:value="selectedVariableData.enum"
                        placeholder="添加可选值"
                      />
                    </n-form-item>
                  </n-grid-item>
                </template>

                <template v-if="selectedVariableData.type === 'array'">
                  <n-grid-item :span="2">
                    <div class="form-section-title">数组配置</div>
                  </n-grid-item>
                  <n-grid-item>
                    <n-form-item label="元素类型 (items.type)">
                      <n-select 
                        v-model:value="selectedVariableData.items.type"
                        :options="basicTypeOptions"
                      />
                    </n-form-item>
                  </n-grid-item>
                  <n-grid-item>
                    <n-form-item label="最小项数 (minItems)">
                      <n-input-number v-model:value="selectedVariableData.minItems" :min="0" />
                    </n-form-item>
                  </n-grid-item>
                </template>

                <!-- 数值验证 -->
                <template v-if="selectedVariableData.type === 'integer' || selectedVariableData.type === 'number'">
                  <n-grid-item :span="2">
                    <div class="form-section-title">数值验证</div>
                  </n-grid-item>
                  <n-grid-item>
                    <n-form-item label="最小值 (min)">
                      <n-input-number v-model:value="selectedVariableData.min" />
                    </n-form-item>
                  </n-grid-item>
                  <n-grid-item>
                    <n-form-item label="最大值 (max)">
                      <n-input-number v-model:value="selectedVariableData.max" />
                    </n-form-item>
                  </n-grid-item>
                </template>

                <!-- 字符串验证 -->
                <template v-if="selectedVariableData.type === 'string'">
                  <n-grid-item :span="2">
                    <div class="form-section-title">字符串验证</div>
                  </n-grid-item>
                  <n-grid-item :span="2">
                    <n-form-item label="正则表达式 (pattern)">
                      <n-input v-model:value="selectedVariableData.pattern" placeholder="^[a-zA-Z0-9]+$" />
                    </n-form-item>
                  </n-grid-item>
                </template>

                <!-- 对象数组配置 -->
                <template v-if="selectedVariableData.type === 'object_arr'">
                  <n-grid-item :span="2">
                    <div class="form-section-title">对象数组配置</div>
                  </n-grid-item>
                  <n-grid-item :span="2">
                    <n-form-item label="对象结构定义">
                      <div class="object-structure-hint">
                        <n-text depth="3" style="font-size: 12px;">
                          对象数组中每个对象的结构定义。右键点击下方空白区域添加字段。
                        </n-text>
                      </div>
                      <div class="object-properties-editor">
                        <div v-if="!selectedVariableData.items || !selectedVariableData.items.properties || Object.keys(selectedVariableData.items.properties).length === 0" 
                             class="empty-properties-hint"
                             @contextmenu="onObjectPropertiesContextMenu">
                          暂无字段定义（右键添加）
                        </div>
                        <div v-else class="properties-list">
                          <div v-for="(prop, propName) in selectedVariableData.items.properties" 
                               :key="propName" 
                               class="property-item">
                            <div class="property-header">
                              <span class="property-name">{{ propName }}</span>
                              <n-tag :type="getTypeTagType(prop.type)" size="small">
                                {{ getTypeLabel(prop.type) }}
                              </n-tag>
                              <n-button 
                                size="tiny" 
                                quaternary 
                                type="error"
                                @click="removeObjectProperty(propName)">
                                删除
                              </n-button>
                            </div>
                            <div class="property-description" v-if="prop.description">
                              {{ prop.description }}
                            </div>
                          </div>
                        </div>
                      </div>
                    </n-form-item>
                  </n-grid-item>
                  <n-grid-item>
                    <n-form-item label="最小项数 (minItems)">
                      <n-input-number v-model:value="selectedVariableData.minItems" :min="0" />
                    </n-form-item>
                  </n-grid-item>
                  <n-grid-item>
                    <n-form-item label="最大项数 (maxItems)">
                      <n-input-number v-model:value="selectedVariableData.maxItems" :min="0" />
                    </n-form-item>
                  </n-grid-item>
                </template>

                <!-- UI配置 -->
                <n-grid-item :span="2">
                  <div class="form-section-title">UI配置</div>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="面板显示 (ui.panel)">
                    <n-switch v-model:value="selectedVariableData.ui.panel">
                      <template #checked>显示</template>
                      <template #unchecked>隐藏</template>
                    </n-switch>
                  </n-form-item>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="排序权重 (ui.order)">
                    <n-input-number v-model:value="selectedVariableData.ui.order" :step="10" />
                  </n-form-item>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="分组标题 (ui.group)">
                    <n-input v-model:value="selectedVariableData.ui.group" placeholder="基础信息" />
                  </n-form-item>
                </n-grid-item>
                
                <n-grid-item>
                  <n-form-item label="组件类型 (ui.component)">
                    <n-select 
                      v-model:value="selectedVariableData.ui.component"
                      :options="componentOptions"
                    />
                  </n-form-item>
                </n-grid-item>

              </n-grid>
            </n-form>
          </div>

          <!-- 空状态 -->
          <div v-else class="empty-selection">
            <n-empty description="请从左侧选择一个变量进行编辑">
              <template #icon>
                <n-icon><SettingsOutline /></n-icon>
              </template>
            </n-empty>
          </div>
        </div>

        <!-- 中右分割线 -->
        <div class="panel-resizer right-resizer" @mousedown="startRightResize"></div>

        <!-- 右侧：预览 -->
        <div class="right-panel" :style="{ width: rightPanelWidth + 'px' }">
          <div class="panel-header">
            <div class="panel-title-row">
              <div class="panel-title">Schema预览</div>
              <!-- 格式切换 -->
              <n-button-group size="small">
                <n-button 
                  :type="previewFormat === 'json' ? 'primary' : 'default'"
                  @click="previewFormat = 'json'"
                >
                  JSON
                </n-button>
                <n-button 
                  :type="previewFormat === 'yaml' ? 'primary' : 'default'"
                  @click="previewFormat = 'yaml'"
                >
                  YAML
                </n-button>
              </n-button-group>
            </div>
            
            <div class="panel-actions">
              <!-- 操作按钮 -->
              <n-space size="small">
                <n-button size="small" @click="importFromFile" quaternary>
                  <template #icon>
                    <n-icon><CloudUploadOutline /></n-icon>
                  </template>
                  导入
                </n-button>
                
                <n-button size="small" @click="showExportDialog" quaternary>
                  <template #icon>
                    <n-icon><DownloadOutline /></n-icon>
                  </template>
                  导出
                </n-button>
                
                <n-button size="small" @click="copyToClipboard" quaternary>
                  <template #icon>
                    <n-icon><CopyOutline /></n-icon>
                  </template>
                  复制
                </n-button>
                
                <n-button size="small" @click="formatContent" quaternary>
                  <template #icon>
                    <n-icon><RefreshOutline /></n-icon>
                  </template>
                  格式化
                </n-button>
                
                <n-button size="small" @click="syncToTree" quaternary>
                  <template #icon>
                    <n-icon><SyncOutline /></n-icon>
                  </template>
                  同步到树
                </n-button>
              </n-space>
            </div>
          </div>
          <div class="preview-content">
            <div ref="schemaEditorRef" class="code-preview"></div>
          </div>
        </div>
      </div>
    </n-drawer-content>
  </n-drawer>

  <!-- 导出格式选择对话框 -->
  <n-modal v-model:show="showExportModalFlag" :mask-closable="false">
    <n-card style="width: 400px" title="选择导出格式" :bordered="false" size="huge">
      <template #header-extra>
        <n-button quaternary circle @click="showExportModalFlag = false">
          <template #icon>
            <n-icon><CloseOutline /></n-icon>
          </template>
        </n-button>
      </template>
      
      <n-space vertical>
        <div>请选择要导出的文件格式：</div>
        <n-button-group size="large" vertical>
          <n-button @click="exportFile('json')" style="justify-content: flex-start;">
            <template #icon>
              <n-icon><CodeOutline /></n-icon>
            </template>
            JSON 格式 (.json)
          </n-button>
          <n-button @click="exportFile('yaml')" style="justify-content: flex-start;">
            <template #icon>
              <n-icon><DocumentOutline /></n-icon>
            </template>
            YAML 格式 (.yaml)
          </n-button>
        </n-button-group>
      </n-space>
    </n-card>
  </n-modal>

  <!-- 导入确认对话框 -->
  <n-modal v-model:show="showImportConfirmModal" :mask-closable="false">
    <n-card style="width: 500px" title="导入确认" :bordered="false" size="huge">
      <template #header-extra>
        <n-button quaternary circle @click="showImportConfirmModal = false">
          <template #icon>
            <n-icon><CloseOutline /></n-icon>
          </template>
        </n-button>
      </template>
      
      <div>
        <p>确定要导入此schema吗？这将替换当前的所有变量定义。</p>
      </div>
      
      <template #footer>
        <div style="display: flex; justify-content: flex-end; gap: 12px;">
          <n-button @click="showImportConfirmModal = false">取消</n-button>
          <n-button type="primary" @click="confirmImport">确定导入</n-button>
        </div>
      </template>
    </n-card>
  </n-modal>

  <!-- 同步确认对话框 -->
  <n-modal v-model:show="showSyncConfirmModal" :mask-closable="false">
    <n-card style="width: 500px" title="同步确认" :bordered="false" size="huge">
      <template #header-extra>
        <n-button quaternary circle @click="showSyncConfirmModal = false">
          <template #icon>
            <n-icon><CloseOutline /></n-icon>
          </template>
        </n-button>
      </template>
      
      <div>
        <p>确定要将预览内容同步到变量树吗？这将替换当前的变量定义。</p>
      </div>
      
      <template #footer>
        <div style="display: flex; justify-content: flex-end; gap: 12px;">
          <n-button @click="showSyncConfirmModal = false">取消</n-button>
          <n-button type="warning" @click="confirmSync">确定同步</n-button>
        </div>
      </template>
    </n-card>
  </n-modal>

  <!-- 测试数据模态框 -->
  <n-modal v-model:show="showTestDataModalFlag" :mask-closable="false">
    <n-card style="width: 900px; max-height: 700px;" title="生成测试数据" :bordered="false" size="huge">
      <template #header-extra>
        <n-button quaternary circle @click="showTestDataModalFlag = false">
          <template #icon>
            <n-icon><CloseOutline /></n-icon>
          </template>
        </n-button>
      </template>
      
      <n-space vertical>
        <div style="margin-bottom: 12px;">
          <n-text depth="3">
            基于当前变量定义自动生成的测试数据，你可以直接编辑这些数据用于测试模板。
          </n-text>
        </div>
        
        <div style="display: flex; justify-content: flex-end; gap: 8px; margin-bottom: 12px;">
          <n-button size="small" @click="regenerateTestData" quaternary>
            <template #icon>
              <n-icon><RefreshOutline /></n-icon>
            </template>
            重新生成
          </n-button>
          <n-button size="small" @click="copyTestData" quaternary>
            <template #icon>
              <n-icon><CopyOutline /></n-icon>
            </template>
            复制数据
          </n-button>
          <n-button size="small" @click="saveTestDataToStorage" type="primary">
            <template #icon>
              <n-icon><SaveOutline /></n-icon>
            </template>
            保存数据
          </n-button>
        </div>
        
        <div class="test-data-editor" ref="testDataEditor" style="min-height: 400px; border: 1px solid #e0e0e6; border-radius: 6px;"></div>
      </n-space>
      
      <template #footer>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <n-text depth="3" style="font-size: 12px;">
            数据将保存到 localStorage，键名：template_studio_{{ props.templateId }}_testdata
          </n-text>
          <div style="display: flex; gap: 12px;">
            <n-button @click="saveTestDataToStorage" type="primary">
              <template #icon>
                <n-icon><SaveOutline /></n-icon>
              </template>
              保存并应用
            </n-button>
            <n-button @click="showTestDataModalFlag = false">关闭</n-button>
          </div>
        </div>
      </template>
    </n-card>
  </n-modal>

  <!-- 变量分析模态框 -->
  <n-modal v-model:show="showVariableAnalysisModalFlag" :mask-closable="false">
    <n-card style="width: 1000px; max-height: 80vh;" title="模板变量分析" :bordered="false" size="huge">
      <template #header-extra>
        <n-button quaternary circle @click="showVariableAnalysisModalFlag = false">
          <template #icon>
            <n-icon><CloseOutline /></n-icon>
          </template>
        </n-button>
      </template>
      
      <n-spin :show="analyzingVariables">
        <div v-if="analysisResult">
          <!-- 调试信息 -->
          <n-alert v-if="true" type="info" style="margin-bottom: 16px;" :show-icon="false">
            调试信息: 缺失变量{{ analysisResult.missingVariables?.length || 0 }}个, 
            检测变量{{ analysisResult.detectedVariables?.length || 0 }}个, 
            未使用变量{{ analysisResult.unusedVariables?.length || 0 }}个
          </n-alert>
          <n-tabs type="line" animated>
            <!-- 缺失变量 -->
            <n-tab-pane name="missing" :tab="`缺失变量 (${analysisResult.missingVariables?.length || 0})`">
              <div v-if="analysisResult.missingVariables && analysisResult.missingVariables.length > 0">
                <div style="margin-bottom: 16px;">
                  <n-text depth="3">模板中使用但未定义的变量：</n-text>
                  <n-button 
                    size="small" 
                    type="primary" 
                    style="margin-left: 12px;"
                    @click="addAllDetectedVariables(analysisResult.missingVariables)"
                  >
                    添加全部到变量树
                  </n-button>
                </div>
                <n-data-table
                  :columns="missingVariablesColumns"
                  :data="analysisResult.missingVariables"
                  :pagination="false"
                  size="small"
                />
              </div>
              <n-empty v-else description="没有缺失的变量" />
            </n-tab-pane>
            
            <!-- 检测到的变量 -->
            <n-tab-pane name="detected" :tab="`检测变量 (${analysisResult.detectedVariables?.length || 0})`">
              <div v-if="analysisResult.detectedVariables && analysisResult.detectedVariables.length > 0">
                <div style="margin-bottom: 16px;">
                  <n-text depth="3">模板中检测到的所有变量：</n-text>
                </div>
                <n-data-table
                  :columns="detectedVariablesColumns"
                  :data="analysisResult.detectedVariables"
                  :pagination="false"
                  size="small"
                />
              </div>
              <n-empty v-else description="没有检测到变量" />
            </n-tab-pane>
            
            <!-- 冲突变量 -->
            <n-tab-pane name="conflict" :tab="`类型冲突 (${analysisResult.conflictVariables?.length || 0})`">
              <div v-if="analysisResult.conflictVariables?.length > 0">
                <div style="margin-bottom: 16px;">
                  <n-text depth="3">已定义但类型可能不匹配的变量：</n-text>
                </div>
                <n-data-table
                  :columns="[
                    { title: '变量名', key: 'name', width: 120 },
                    { title: '推测类型', key: 'type', width: 100 },
                    { title: '出现文件', key: 'files', render: (row) => row.files?.join(', ') },
                    { title: '冲突说明', key: 'suggestions' }
                  ]"
                  :data="analysisResult.conflictVariables"
                  :pagination="false"
                  size="small"
                />
              </div>
              <n-empty v-else description="没有类型冲突" />
            </n-tab-pane>
            
            <!-- 未使用变量 -->
            <n-tab-pane name="unused" :tab="`未使用 (${analysisResult.unusedVariables?.length || 0})`">
              <div v-if="analysisResult.unusedVariables && analysisResult.unusedVariables.length > 0">
                <div style="margin-bottom: 16px;">
                  <n-text depth="3">已定义但模板中未使用的变量：</n-text>
                  <n-button 
                    size="small" 
                    type="warning" 
                    style="margin-left: 12px;"
                    @click="deleteAllUnusedVariables(analysisResult.unusedVariables)"
                  >
                    删除全部未使用变量
                  </n-button>
                </div>
                <n-data-table
                  :columns="unusedVariablesColumns"
                  :data="analysisResult.unusedVariables.map(name => ({ name }))"
                  :pagination="false"
                  size="small"
                />
              </div>
              <n-empty v-else description="所有变量都已使用" />
            </n-tab-pane>
            
            <!-- 统计信息 -->
            <n-tab-pane name="stats" tab="统计信息">
              <n-descriptions label-placement="left" :column="2">
                <n-descriptions-item label="总变量数">
                  <n-tag type="info" size="small">{{ analysisResult.totalVariableCount || 0 }}</n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="分析文件数">
                  <n-tag type="default" size="small">{{ analysisResult.analyzedFileCount || 0 }}</n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="缺失变量">
                  <n-tag :type="(analysisResult.missingVariables?.length || 0) > 0 ? 'warning' : 'success'" size="small">
                    {{ analysisResult.missingVariables?.length || 0 }}
                  </n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="类型冲突">
                  <n-tag :type="(analysisResult.conflictVariables?.length || 0) > 0 ? 'error' : 'success'" size="small">
                    {{ analysisResult.conflictVariables?.length || 0 }}
                  </n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="未使用变量">
                  <n-tag :type="(analysisResult.unusedVariables?.length || 0) > 0 ? 'warning' : 'success'" size="small">
                    {{ analysisResult.unusedVariables?.length || 0 }}
                  </n-tag>
                </n-descriptions-item>
              </n-descriptions>
              
              <!-- 快速操作 -->
              <div style="margin-top: 24px;">
                <n-text depth="3" style="margin-bottom: 12px; display: block;">快速操作：</n-text>
                <n-space>
                  <n-button 
                    size="small" 
                    type="primary" 
                    :disabled="!analysisResult.missingVariables || analysisResult.missingVariables.length === 0"
                    @click="addAllDetectedVariables(analysisResult.missingVariables)"
                  >
                    添加所有缺失变量 ({{ analysisResult.missingVariables?.length || 0 }})
                  </n-button>
                  <n-button 
                    size="small" 
                    type="warning" 
                    :disabled="!analysisResult.unusedVariables || analysisResult.unusedVariables.length === 0"
                    @click="cleanupUnusedVariables"
                  >
                    删除未使用变量 ({{ analysisResult.unusedVariables?.length || 0 }})
                  </n-button>
                  <n-button 
                    size="small" 
                    type="info" 
                    @click="optimizeAllVariableTypes"
                    :disabled="!analysisResult.detectedVariables || analysisResult.detectedVariables.length === 0"
                  >
                    优化所有变量类型
                  </n-button>
                </n-space>
              </div>
            </n-tab-pane>
          </n-tabs>
        </div>
        
        <template #description>
          正在分析模板变量...
        </template>
      </n-spin>
      
      <template #footer>
        <div style="display: flex; justify-content: flex-end; gap: 12px;">
          <n-button @click="showVariableAnalysisModal" :loading="analyzingVariables">重新分析</n-button>
          <n-button @click="showVariableAnalysisModalFlag = false">关闭</n-button>
        </div>
      </template>
    </n-card>
  </n-modal>

  <!-- 隐藏的文件输入框 -->
  <input
    ref="fileInputRef"
    type="file"
    accept=".json,.yaml,.yml"
    style="display: none"
    @change="handleFileImport"
  />
</template>

<script setup>
import { ref, computed, watch, h, nextTick, onMounted, onUnmounted, markRaw, toRaw } from 'vue'
import { 
  NDrawer, NDrawerContent, NButton, NIcon, NSpace, NForm, NFormItem, 
  NGrid, NGridItem, NInput, NSelect, NSwitch, NTabs, NTabPane, NInputNumber,
  NTree, NTag, NDropdown, NDynamicTags, NEmpty, NButtonGroup, NModal, NCard,
  NText, NSpin, NDataTable, NList, NListItem, NDescriptions, NDescriptionsItem, 
  NAlert, useMessage 
} from 'naive-ui'
import { 
  SaveOutline, CloseOutline, DocumentOutline, AddOutline, TrashOutline,
  EllipsisHorizontalOutline, SettingsOutline, CodeSlashOutline,
  TextOutline, Calculator, ToggleOutline, ListOutline, ArchiveOutline,
  KeyOutline, LockClosedOutline, Folder, FolderOpenOutline, 
  EllipsisVerticalOutline, CheckboxOutline, CodeOutline, ChevronForward,
  CreateOutline, DownloadOutline, CloudUploadOutline, CopyOutline, 
  RefreshOutline, SyncOutline, FlaskOutline, SearchOutline
} from '@vicons/ionicons5'
import { EditorView, basicSetup } from 'codemirror'
import { EditorState } from '@codemirror/state'
import { json } from '@codemirror/lang-json'
import { yaml } from '@codemirror/lang-yaml'
import * as YAML from 'js-yaml'
import { getTemplateExpose, setTemplateExpose } from '@/api/templateExpose'
import { analyzeTemplateVariables } from '@/api/templates'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  templateId: {
    type: [String, Number],
    required: true
  },
  templateVariables: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:show', 'variables-saved', 'test-data-updated'])

const message = useMessage()

// 抽屉显示状态
const visible = computed({
  get: () => props.show,
  set: (value) => emit('update:show', value)
})

// 预览面板操作相关
const showExportModalFlag = ref(false)
const showImportConfirmModal = ref(false)
const showSyncConfirmModal = ref(false)
const fileInputRef = ref(null)
const pendingImportData = ref(null)
const pendingSyncData = ref(null)

// 数据状态
const saving = ref(false)
const varsSchema = ref({}) // 变量Schema数据
const selectedKeys = ref([])
const expandedKeys = ref([])
const selectedVariableData = ref(null)
const previewFormat = ref('json')
// 预览面板显示Schema结构定义
const schemaEditorRef = ref(null)
let schemaEditor = null

// 右键菜单状态
const showContextMenuFlag = ref(false)
const contextMenuX = ref(0)
const contextMenuY = ref(0)
const contextMenuOptions = ref([])
const contextMenuTarget = ref(null)

// 编辑状态管理（参考模板文件树）
const editingNode = ref(null)
const renamingNode = ref(null)
const newVariableName = ref('')
const addVariableType = ref('string')

// 拖拽相关状态
const drawerHeight = ref(Math.floor(window.innerHeight * 0.67)) // 默认三分之二高度
const isResizing = ref(false)

// 测试数据相关状态
const showTestDataModalFlag = ref(false)
const testData = ref({})
const testDataEditor = ref(null)

// 变量分析相关状态
const showVariableAnalysisModalFlag = ref(false)
const analysisResult = ref(null)
const analyzingVariables = ref(false)
const minHeight = 300
const maxHeight = window.innerHeight - 100

// 变量分析表格列定义
const missingVariablesColumns = [
  { title: '变量名', key: 'name', width: 120 },
  { 
    title: '推测类型', 
    key: 'type', 
    width: 100,
    render(row) {
      return h(NTag, { 
        type: getTypeTagType(row.type),
        size: 'small'
      }, { default: () => row.type })
    }
  },
  { 
    title: '出现文件', 
    key: 'files', 
    render(row) {
      if (!row.files || row.files.length === 0) return '-'
      return row.files.slice(0, 3).join(', ') + (row.files.length > 3 ? `等${row.files.length}个文件` : '')
    }
  },
  { 
    title: '使用上下文', 
    key: 'contexts', 
    render(row) {
      if (!row.contexts || row.contexts.length === 0) return '-'
      return row.contexts.slice(0, 2).join(', ') + (row.contexts.length > 2 ? '...' : '')
    }
  },
  { title: '建议', key: 'suggestions', width: 200 },
  { 
    title: '操作', 
    key: 'actions',
    width: 80,
    render(row) {
      return h(NButton, { 
        size: 'small',
        type: 'primary',
        onClick: () => addDetectedVariable(row) 
      }, { default: () => '添加' })
    }
  }
]

const detectedVariablesColumns = [
  { title: '变量名', key: 'name', width: 120 },
  { 
    title: '推测类型', 
    key: 'type', 
    width: 100,
    render(row) {
      return h(NTag, { 
        type: getTypeTagType(row.type),
        size: 'small'
      }, { default: () => row.type })
    }
  },
  { 
    title: '出现文件', 
    key: 'files', 
    render(row) {
      if (!row.files || row.files.length === 0) return '-'
      return row.files.slice(0, 3).join(', ') + (row.files.length > 3 ? `等${row.files.length}个文件` : '')
    }
  },
  { 
    title: '使用上下文', 
    key: 'contexts', 
    render(row) {
      if (!row.contexts || row.contexts.length === 0) return '-'
      return row.contexts.slice(0, 2).join(', ') + (row.contexts.length > 2 ? '...' : '')
    }
  },
  { title: '建议', key: 'suggestions', width: 200 }
]

const unusedVariablesColumns = [
  { title: '变量名', key: 'name', width: 200 },
  { 
    title: '当前定义',
    key: 'definition',
    render(row) {
      const currentVar = varsSchema.value[row.name]
      if (!currentVar) return '-'
      return h(NTag, { 
        type: getTypeTagType(currentVar.type),
        size: 'small'
      }, { default: () => currentVar.type })
    }
  },
  { 
    title: '描述',
    key: 'description',
    render(row) {
      const currentVar = varsSchema.value[row.name]
      return currentVar?.description || '-'
    }
  },
  { 
    title: '操作', 
    key: 'actions',
    width: 80,
    render(row) {
      return h(NButton, { 
        size: 'small',
        type: 'error',
        onClick: () => deleteUnusedVariable(row.name) 
      }, { default: () => '删除' })
    }
  }
]

// 面板宽度状态
const leftPanelWidth = ref(250)
const centerPanelWidth = ref(0) // 将通过computed计算
const rightPanelWidth = ref(0) // 将通过computed计算
const isLeftResizing = ref(false)
const isRightResizing = ref(false)
const totalLayoutWidth = ref(1200) // 默认总宽度

// 计算中间和右侧面板宽度
const computedCenterWidth = computed(() => {
  const remaining = totalLayoutWidth.value - leftPanelWidth.value - rightPanelWidth.value - 32 - 16 // 减去padding和间距
  return Math.max(300, remaining) // 最小宽度300px
})

const computedRightWidth = computed(() => {
  return Math.max(300, Math.floor(totalLayoutWidth.value / 3)) // 默认1/3宽度，最小300px
})

// 监听计算值变化
watch(computedCenterWidth, (newWidth) => {
  centerPanelWidth.value = newWidth
})

watch(computedRightWidth, (newWidth) => {
  rightPanelWidth.value = newWidth
}, { immediate: true })

// 节点属性配置（参考TemplateFileTree）
const nodeProps = ({ option }) => {
  return {
    onClick(e) {
      console.log('节点被点击:', option.key)
      // 手动设置选中状态
      selectedKeys.value = [option.key]
      // 触发选择处理
      onSelectVariable([option.key])
    },
    onContextmenu(e) {
      e.preventDefault()
      e.stopPropagation()
      
      // 设置下拉菜单选项
      contextMenuOptions.value = getNodeMenuOptions(option)
      contextMenuTarget.value = option
      contextMenuX.value = e.clientX
      contextMenuY.value = e.clientY
      showContextMenuFlag.value = true
    }
  }
}

// 隐藏右键菜单
const hideContextMenu = () => {
  showContextMenuFlag.value = false
  contextMenuTarget.value = null
}

// 处理右键菜单操作
const handleContextMenuAction = (key) => {
  if (contextMenuTarget.value) {
    handleNodeAction(key, contextMenuTarget.value)
  } else {
    // 空白区域的右键菜单操作
    handleEmptyAreaAction(key)
  }
  hideContextMenu()
}

// 空白区域右键菜单处理
const handleEmptyAreaAction = (key) => {
  switch (key) {
    case 'add-string':
      addRootVariable() // 添加普通变量（默认string类型）
      break
    case 'add-object':
      addRootObjectVariable() // 添加对象变量
      break
  }
}


// 树区域右键菜单处理
const onTreeAreaContextMenu = (event) => {
  // 检查是否点击在树节点上
  if (event.target.closest('.n-tree-node')) return
  
  event.preventDefault()
  event.stopPropagation()
  
  // 设置空白区域的右键菜单选项
  contextMenuOptions.value = [
    {
      label: '添加变量',
      key: 'add-string',
      icon: () => h(NIcon, null, { default: () => h(AddOutline) })
    },
    {
      label: '添加对象变量',
      key: 'add-object',
      icon: () => h(NIcon, null, { default: () => h(ArchiveOutline) })
    }
  ]
  
  contextMenuTarget.value = null // 标记为空白区域
  contextMenuX.value = event.clientX
  contextMenuY.value = event.clientY
  showContextMenuFlag.value = true
  
}

// 数据类型选项
const typeOptions = [
  { label: '字符串 (string)', value: 'string' },
  { label: '整数 (integer)', value: 'integer' },
  { label: '数字 (number)', value: 'number' },
  { label: '布尔值 (boolean)', value: 'boolean' },
  { label: '数组 (array)', value: 'array' },
  { label: '对象 (object)', value: 'object' },
  { label: '对象数组 (object_arr)', value: 'object_arr' },
  { label: '枚举 (enum)', value: 'enum' },
  { label: '密码 (secret)', value: 'secret' }
]

const basicTypeOptions = [
  { label: '字符串', value: 'string' },
  { label: '整数', value: 'integer' },
  { label: '数字', value: 'number' },
  { label: '布尔值', value: 'boolean' }
]

const componentOptions = [
  { label: '输入框 (input)', value: 'input' },
  { label: '下拉选择 (select)', value: 'select' },
  { label: '开关 (switch)', value: 'switch' },
  { label: '文本域 (textarea)', value: 'textarea' }
]


// 转换为变量树数据 - 支持编辑节点
const variableTreeData = computed(() => {
  try {
    const treeData = convertToTreeData(varsSchema.value)
    
    // 如果有编辑节点，插入到正确位置
    if (editingNode.value) {
      insertEditingNodeToTree(treeData, editingNode.value)
    }
    
    // console.log('生成的树数据:', treeData) // 调试用，已注释
    return treeData
  } catch (error) {
    console.error('生成变量树数据时出错:', error)
    return []
  }
})

// 将编辑节点插入到树的正确位置
const insertEditingNodeToTree = (treeData, editingNode) => {
  // 创建编辑节点
  const createEditingNode = () => {
    const nodeType = editingNode.type || 'string'
    return {
      key: editingNode.key,
      title: newVariableName.value || '',
      type: nodeType,
      path: editingNode.key,
      data: null,
      isLeaf: nodeType !== 'object',
      isEditing: true,
      children: [],
      prefix: () => h(NIcon, { class: `var-icon var-${nodeType}` }, {
        default: () => h(getVariableIconComponent(nodeType, false))
      })
    }
  }
  
  // 如果是根节点编辑
  if (editingNode.isRoot || !editingNode.parentPath) {
    treeData.unshift(createEditingNode())
    return true
  }
  
  // 否则查找父节点并插入
  const findAndInsert = (nodes) => {
    for (const node of nodes) {
      if (node.path === editingNode.parentPath) {
        if (!node.children) node.children = []
        node.children.unshift(createEditingNode())
        return true
      }
      if (node.children && node.children.length > 0 && findAndInsert(node.children)) {
        return true
      }
    }
    return false
  }
  
  const result = findAndInsert(treeData)
  if (!result) {
    console.error('Failed to find parent node for editing node:', editingNode.parentPath)
  }
  return result
}

// 转换变量Schema为树形数据
const convertToTreeData = (schema, parentPath = '', visited = new Set(), depth = 0) => {
  const treeData = []
  
  // 防止过深的递归
  if (depth > 10) {
    console.warn('递归深度超过限制，停止处理:', parentPath)
    return treeData
  }
  
  if (!schema || typeof schema !== 'object') {
    return treeData
  }
  
  try {
    Object.entries(schema).forEach(([key, value]) => {
      const currentPath = parentPath ? `${parentPath}.${key}` : key
      
      // 检测循环引用
      if (visited.has(currentPath)) {
        console.warn('检测到循环引用，跳过节点:', currentPath)
        return
      }
      
      // 额外检查：避免同名属性造成的递归路径问题
      const pathParts = currentPath.split('.')
      if (pathParts.length > 2) {
        // 检查是否存在重复的相邻路径段
        for (let i = 1; i < pathParts.length - 1; i++) {
          if (pathParts[i] === pathParts[i + 1]) {
            console.warn('检测到重复的相邻路径段，跳过节点:', currentPath)
            return
          }
        }
      }
      
      // 检查值是否有效
      if (!value || typeof value !== 'object') {
        console.warn('无效的变量值，跳过:', key, value)
        return
      }
      
      const varType = value.type || 'string'
      
      // 检查对象类型和对象数组类型是否有子变量
      const hasChildren = (varType === 'object' && value.properties && 
                          typeof value.properties === 'object' && 
                          Object.keys(value.properties).length > 0) ||
                         (varType === 'object_arr' && value.items && value.items.properties &&
                          typeof value.items.properties === 'object' &&
                          Object.keys(value.items.properties).length > 0)
      
      const node = {
        key: currentPath,
        title: key,
        type: varType,
        path: currentPath,
        data: value,
        isLeaf: varType !== 'object' && varType !== 'object_arr',
        isEditing: Boolean(renamingNode.value && renamingNode.value.path === currentPath),
        children: [],
        prefix: () => h(NIcon, { class: `var-icon var-${varType}` }, {
          default: () => h(getVariableIconComponent(varType, hasChildren))
        })
      }
      
      // 如果是对象类型或对象数组类型且有properties，递归生成子节点
      if (hasChildren) {
        try {
          const newVisited = new Set(visited)
          newVisited.add(currentPath)
          
          // 根据类型选择正确的properties路径
          let childProperties
          if (varType === 'object') {
            childProperties = value.properties
          } else if (varType === 'object_arr') {
            childProperties = value.items.properties
          }
          
          node.children = convertToTreeData(childProperties, currentPath, newVisited, depth + 1)
        } catch (childError) {
          console.error('处理子节点时出错:', currentPath, childError)
          node.children = []
        }
      }
      
      treeData.push(node)
    })
  } catch (error) {
    console.error('convertToTreeData 处理时出错:', error)
  }
  
  return treeData
}

// 获取变量图标组件（参考模板资源树的实现）
const getVariableIconComponent = (type, hasChildren = false) => {
  // 对象类型根据是否有子变量显示不同图标
  if (type === 'object') {
    return hasChildren ? FolderOpenOutline : Folder
  }
  
  // 对象数组类型使用特殊图标
  if (type === 'object_arr') {
    return hasChildren ? FolderOpenOutline : ListOutline
  }
  
  // 其他类型使用标准图标
  const iconMap = {
    string: TextOutline,        // 📝 文本图标
    integer: Calculator,        // 🔢 计算器图标
    number: Calculator,         // 🔢 计算器图标 
    boolean: CheckboxOutline,   // ☑️ 复选框图标
    array: ListOutline,         // 📋 列表图标
    enum: EllipsisVerticalOutline, // ⋮ 选择图标
    secret: LockClosedOutline   // 🔒 锁图标
  }
  return iconMap[type] || TextOutline
}

// 渲染展开/收起图标（参考模板文件树）
const renderSwitcherIcon = () => h(NIcon, null, { default: () => h(ChevronForward) })

// 渲染标签（支持内联编辑）
const renderLabel = ({ option }) => {
  if (option.isEditing === true) {
    const isRenaming = renamingNode.value && String(option.key) === String(renamingNode.value.key)
    const placeholder = isRenaming ? '重命名变量' : '输入变量名'
    
    return h('input', {
      class: 'vscode-tree-input',
      value: newVariableName.value,
      autofocus: true,
      placeholder: placeholder,
      onInput: e => (newVariableName.value = e.target.value),
      onKeydown: e => {
        if (e.key === 'Enter') {
          e.preventDefault()
          confirmAddVariable()
        }
        if (e.key === 'Escape') {
          e.preventDefault()
          cancelAddVariable()
        }
      },
      onMousedown: e => {
        // 防止输入框点击时失去焦点
        e.stopPropagation()
      },
      onBlur: (e) => {
        // 延迟执行，避免与点击事件冲突
        setTimeout(() => {
          // 检查焦点是否还在输入框内，如果是则不处理
          if (document.activeElement && document.activeElement.classList.contains('vscode-tree-input')) {
            return
          }
          
          if (editingNode.value || renamingNode.value) {
            confirmAddVariable()
          }
        }, 150)
      }
    })
  }
  
  return option.title
}

// 安全地访问嵌套变量的辅助函数
const getVariableByPath = (path) => {
  if (!path) return null
  
  const pathParts = path.split('.')
  let current = varsSchema.value
  
  for (let i = 0; i < pathParts.length; i++) {
    const part = pathParts[i]
    
    if (i === 0) {
      // 第一层直接访问
      if (current[part]) {
        current = current[part]
      } else {
        return null
      }
    } else {
      // 后续层需要通过properties或items.properties访问
      if (current.properties && current.properties[part]) {
        current = current.properties[part]
      } else if (current.items && current.items.properties && current.items.properties[part]) {
        // 处理object_arr类型的子变量
        current = current.items.properties[part]
      } else {
        return null
      }
    }
  }
  
  return current
}

// 安全地设置嵌套变量的辅助函数
const setVariableByPath = (path, value) => {
  if (!path) return false
  
  const pathParts = path.split('.')
  const varName = pathParts[pathParts.length - 1]
  
  if (pathParts.length === 1) {
    // 设置根级变量
    varsSchema.value[varName] = value
    return true
  }
  
  // 获取父级变量
  const parentPath = pathParts.slice(0, -1).join('.')
  const parent = getVariableByPath(parentPath)
  
  if (parent && parent.type === 'object') {
    if (!parent.properties) {
      parent.properties = {}
    }
    parent.properties[varName] = value
    return true
  }
  
  return false
}

// 安全地删除嵌套变量的辅助函数
const deleteVariableByPath = (path) => {
  if (!path) return false
  
  const pathParts = path.split('.')
  const varName = pathParts[pathParts.length - 1]
  
  if (pathParts.length === 1) {
    // 删除根级变量
    delete varsSchema.value[varName]
    return true
  }
  
  // 获取父级变量
  const parentPath = pathParts.slice(0, -1).join('.')
  const parent = getVariableByPath(parentPath)
  
  if (parent && parent.properties && parent.properties[varName]) {
    delete parent.properties[varName]
    return true
  }
  
  return false
}

// 获取变量容器（用于添加同级变量）
const getVariableContainer = (path) => {
  if (!path) return varsSchema.value
  
  const pathParts = path.split('.')
  
  if (pathParts.length === 1) {
    // 根级变量的容器就是根对象
    return varsSchema.value
  }
  
  // 获取父级变量
  const parentPath = pathParts.slice(0, -1).join('.')
  const parent = getVariableByPath(parentPath)
  
  if (parent && parent.type === 'object') {
    if (!parent.properties) {
      parent.properties = {}
    }
    return parent.properties
  }
  
  return null
}

// 获取类型标签样式
const getTypeTagType = (type) => {
  const tagMap = {
    string: 'info',
    integer: 'success',
    number: 'success', 
    boolean: 'warning',
    array: 'error',
    object: 'default',
    object_arr: 'primary',
    enum: 'info',
    secret: 'error',
    unknown: 'warning'
  }
  return tagMap[type] || 'default'
}

// 获取类型标签
const getTypeLabel = (type) => {
  const typeMap = {
    string: '字符串',
    integer: '整数',
    number: '数字',
    boolean: '布尔值',
    array: '数组',
    object: '对象',
    object_arr: '对象数组',
    enum: '枚举',
    secret: '密码'
  }
  return typeMap[type] || type
}

// 获取输入组件
const getInputComponent = (type) => {
  const componentMap = {
    string: 'n-input',
    integer: 'n-input-number',
    number: 'n-input-number',
    boolean: 'n-switch',
    secret: 'n-input'
  }
  return componentMap[type] || 'n-input'
}

// 获取变量模板字符串
const getVariableTemplate = (varName) => {
  return `{{.${varName || 'varName'}}}`
}

// 创建默认变量配置
const createDefaultVariable = (name = '', type = 'string') => {
  const variable = {
    type,
    title: name,
    description: '',
    required: false,
    default: getDefaultValue(type),
    insertText: getVariableTemplate(name),
    ui: {
      panel: false,
      order: 10,
      group: '基础信息',
      component: getDefaultComponent(type)
    }
  }

  // 只为特定类型添加相应的字段
  if (type === 'enum') {
    variable.enum = []
  }
  
  if (type === 'array') {
    variable.items = { type: 'string' }
  }
  
  if (type === 'object') {
    variable.properties = {}
  }
  
  if (type === 'object_arr') {
    variable.items = {
      type: 'object',
      properties: {}
    }
  }
  
  if (['integer', 'number'].includes(type)) {
    // 为数值类型预留min/max字段，但不设置默认值
  }
  
  if (type === 'string') {
    // 为字符串类型预留pattern字段，但不设置默认值
  }

  return variable
}

// 获取类型默认值
const getDefaultValue = (type) => {
  const defaults = {
    string: '',
    integer: 0,
    number: 0.0,
    boolean: false,
    array: [],
    object: {},
    object_arr: [],
    enum: '',
    secret: ''
  }
  return defaults[type]
}

// 获取类型默认组件
const getDefaultComponent = (type) => {
  const components = {
    string: 'input',
    integer: 'input',
    number: 'input',
    boolean: 'switch',
    array: 'input',
    object: 'input',
    object_arr: 'object_array_editor',
    enum: 'select',
    secret: 'input'
  }
  return components[type] || 'input'
}

// 对象数组属性管理
const onObjectPropertiesContextMenu = (e) => {
  e.preventDefault()
  addObjectProperty()
}

const addObjectProperty = () => {
  const propertyName = window.prompt('请输入属性名称:')
  if (propertyName && propertyName.trim()) {
    const trimmedName = propertyName.trim()
    
    // 确保items对象存在
    if (!selectedVariableData.value.items) {
      selectedVariableData.value.items = { type: 'object', properties: {} }
    }
    if (!selectedVariableData.value.items.properties) {
      selectedVariableData.value.items.properties = {}
    }
    
    selectedVariableData.value.items.properties[trimmedName] = createDefaultVariable(trimmedName, 'string')
  }
}

const removeObjectProperty = (propertyName) => {
  if (selectedVariableData.value.items && selectedVariableData.value.items.properties) {
    delete selectedVariableData.value.items.properties[propertyName]
  }
}

// 树选择处理
const onSelectVariable = (selectedKeys) => {
  try {
    console.log('onSelectVariable 被调用:', selectedKeys)
    
    // 如果正在编辑中且点击的是编辑节点，不处理
    if ((editingNode.value || renamingNode.value) && 
        selectedKeys.length > 0 && 
        (selectedKeys[0].startsWith('__new__') || selectedKeys[0] === editingNode.value?.key)) {
      console.log('跳过编辑节点')
      return
    }
    
    // 如果正在编辑但点击的是其他节点，先取消编辑状态
    if (editingNode.value || renamingNode.value) {
      console.log('取消编辑状态')
      cancelAddVariable()
    }
    
    if (selectedKeys.length > 0) {
      const selectedKey = selectedKeys[0]
      console.log('选中的key:', selectedKey)
      
      // 跳过临时编辑节点
      if (selectedKey.startsWith('__new__')) {
        console.log('跳过临时编辑节点')
        return
      }
      
      const variable = getVariableByPath(selectedKey)
      
      if (variable) {
        // 创建编辑数据的安全副本，避免循环引用
        const safeClone = (obj, visited = new WeakSet()) => {
          if (obj === null || typeof obj !== 'object') return obj
          if (visited.has(obj)) return {} // 跳过循环引用
          
          visited.add(obj)
          const clone = Array.isArray(obj) ? [] : {}
          
          for (const key in obj) {
            if (obj.hasOwnProperty(key)) {
              clone[key] = safeClone(obj[key], visited)
            }
          }
          
          return clone
        }
        
        // 创建安全的编辑数据，排除properties避免循环引用
        const { properties, ...safeVariable } = variable
        
        selectedVariableData.value = {
          path: selectedKey,
          ...safeVariable,
          // 确保必要的字段存在
          insertText: variable.insertText || getVariableTemplate(selectedKey.split('.').pop()),
          // 确保必要的嵌套对象存在
          ui: variable.ui || {
            panel: false,
            order: 10,
            group: '基础信息',
            component: getDefaultComponent(variable.type || 'string')
          },
          items: variable.items || (variable.type === 'array' ? { type: 'string' } : undefined),
          enum: variable.enum || (variable.type === 'enum' ? [] : undefined)
        }
        console.log('设置selectedVariableData成功:', selectedVariableData.value)
      } else {
        selectedVariableData.value = null
      }
    } else {
      selectedVariableData.value = null
    }
  } catch (error) {
    console.error('选择变量时出错:', error)
    selectedVariableData.value = null
  }
}

// 展开键处理
const onExpandKeys = (keys) => {
  expandedKeys.value = keys
}

// 添加根变量 - 恢复内联编辑
const addRootVariable = () => {
  // 创建临时编辑节点
  const tempKey = `__new__${Date.now()}`
  editingNode.value = {
    key: tempKey,
    parentPath: '',
    isRoot: true,
    type: 'string' // 默认string类型
  }
  newVariableName.value = ''
  
  // 强制更新树数据以显示编辑节点
  varsSchema.value = { ...varsSchema.value }
}

// 添加根对象变量 - 直接创建object类型
const addRootObjectVariable = () => {
  // 创建临时编辑节点
  const tempKey = `__new__object_${Date.now()}`
  editingNode.value = {
    key: tempKey,
    parentPath: '',
    isRoot: true,
    type: 'object' // 直接创建object类型
  }
  newVariableName.value = ''
  
  // 强制更新树数据以显示编辑节点
  varsSchema.value = { ...varsSchema.value }
}


// 获取节点菜单选项
const getNodeMenuOptions = (option) => {
  const menuOptions = []
  
  // 对象类型和对象数组类型都可以添加子属性
  if (option.type === 'object' || option.type === 'object_arr') {
    menuOptions.push({
      label: '新增子变量',
      key: 'add-child',
      icon: () => h(NIcon, null, { default: () => h(AddOutline) })
    })
  }
  
  // 重命名变量
  menuOptions.push({
    label: '重命名变量',
    key: 'rename',
    icon: () => h(NIcon, null, { default: () => h(CreateOutline) })
  })
  
  // 复制变量
  menuOptions.push({
    label: '复制变量',
    key: 'copy',
    icon: () => h(NIcon, null, { default: () => h(DocumentOutline) })
  })
  
  // 删除变量
  menuOptions.push({
    label: '删除变量',
    key: 'delete',
    icon: () => h(NIcon, null, { default: () => h(TrashOutline) })
  })
  
  return menuOptions
}

// 处理节点操作
const handleNodeAction = (key, option) => {
  switch (key) {
    case 'add-child':
      addChildVariable(option.path)
      break
    case 'rename':
      renameVariable(option.path)
      break
    case 'copy':
      copyVariable(option.path)
      break
    case 'delete':
      deleteVariable(option.path)
      break
  }
}

// 重命名变量
const renameVariable = (path) => {
  const variable = getVariableByPath(path)
  
  if (!variable) {
    message.error('变量不存在')
    return
  }
  
  // 取消任何现有的编辑状态
  cancelAddVariable()
  
  // 获取当前变量名
  const pathParts = path.split('.')
  const currentName = pathParts[pathParts.length - 1]
  
  // 设置重命名状态
  renamingNode.value = {
    path: path,
    oldName: currentName,
    key: path
  }
  
  // 预填充当前名称
  newVariableName.value = currentName
  
  // 强制更新树数据，确保重命名状态生效
  nextTick(() => {
    varsSchema.value = { ...varsSchema.value }
    
    // 再次确保输入框获得焦点
    nextTick(() => {
      const input = document.querySelector('.vscode-tree-input')
      if (input) {
        input.focus()
        input.select()  // 选中所有文本便于重新输入
      }
    })
  })
}

// 添加子变量 - 恢复内联编辑
const addChildVariable = (parentPath) => {
  // 1. 确保父级变量存在
  const parent = getVariableByPath(parentPath)
  if (!parent) {
    message.error('父级变量不存在')
    return
  }
  
  // 2. 检查父级是否可以包含子变量
  if (parent.type !== 'object' && parent.type !== 'object_arr') {
    message.error('只有对象类型或对象数组类型的变量才能添加子变量')
    return
  }
  
  // 3. 确保有properties对象
  let targetProperties
  if (parent.type === 'object') {
    if (!parent.properties) {
      parent.properties = {}
    }
    targetProperties = parent.properties
  } else if (parent.type === 'object_arr') {
    if (!parent.items) {
      parent.items = { type: 'object', properties: {} }
    }
    if (!parent.items.properties) {
      parent.items.properties = {}
    }
    targetProperties = parent.items.properties
  }
  
  // 4. 创建临时编辑节点
  const tempKey = `__new__${parentPath}_${Date.now()}`
  editingNode.value = {
    key: tempKey,
    parentPath: parentPath,
    isRoot: false
  }
  newVariableName.value = ''
  addVariableType.value = 'string' // 子变量默认为字符串类型
  
  // 5. 展开父级
  if (!expandedKeys.value.includes(parentPath)) {
    expandedKeys.value = [...expandedKeys.value, parentPath]
  }
  
  // 6. 不需要强制更新，Vue会自动检测到数据变化
}


// 复制变量
const copyVariable = (path) => {
  const originalVariable = getVariableByPath(path)
  
  if (!originalVariable) {
    message.error('原始变量不存在')
    return
  }
  
  const pathParts = path.split('.')
  const originalVarName = pathParts[pathParts.length - 1]
  const copyVarName = `${originalVarName}_copy_${Date.now()}`
  
  // 获取变量容器
  const container = getVariableContainer(path)
  
  if (!container) {
    message.error('无法找到变量容器')
    return
  }
  
  // 深度复制变量定义
  const copiedVariable = JSON.parse(JSON.stringify(originalVariable))
  copiedVariable.title = (copiedVariable.title || originalVarName) + ' 复制'
  container[copyVarName] = copiedVariable
  
  // 生成新的路径
  const newPath = pathParts.slice(0, -1).concat([copyVarName]).join('.')
  
  // 强制更新树数据
  varsSchema.value = { ...varsSchema.value }
  
  // 选中新复制的变量
  selectedKeys.value = [newPath]
  
  // 等待DOM更新后再选择变量
  nextTick(() => {
    onSelectVariable([newPath])
  })
  
  message.success(`变量复制成功: ${copyVarName}`)
}


// 删除变量
const deleteVariable = (path) => {
  if (!deleteVariableByPath(path)) {
    message.error('删除变量失败')
    return
  }
  
  // 清除选择（如果删除的是当前选中的变量）
  if (selectedKeys.value.includes(path)) {
    selectedKeys.value = []
    selectedVariableData.value = null
  }
  
  // 强制更新树数据
  varsSchema.value = { ...varsSchema.value }
  
  const pathParts = path.split('.')
  const varName = pathParts[pathParts.length - 1]
  message.success(`已删除变量: ${varName}`)
}

// 类型改变处理
const onTypeChange = (newType) => {
  if (!selectedVariableData.value) return
  
  console.log('类型变更:', selectedVariableData.value.path, '从', selectedVariableData.value.type, '到', newType)
  
  // 重置类型相关的字段
  selectedVariableData.value.default = getDefaultValue(newType)
  selectedVariableData.value.ui.component = getDefaultComponent(newType)
  
  // 更新插入文本（如果当前是默认值，则更新为新的默认值）
  const currentVarName = selectedVariableData.value.path.split('.').pop()
  const expectedInsertText = getVariableTemplate(currentVarName)
  if (!selectedVariableData.value.insertText || selectedVariableData.value.insertText === expectedInsertText) {
    selectedVariableData.value.insertText = expectedInsertText
  }
  
  // 根据类型设置特殊字段
  if (newType === 'enum') {
    selectedVariableData.value.enum = []
  } else {
    delete selectedVariableData.value.enum
  }
  
  if (newType === 'array') {
    selectedVariableData.value.items = { type: 'string' }
    selectedVariableData.value.minItems = undefined
  } else if (newType === 'object_arr') {
    selectedVariableData.value.items = { type: 'object', properties: {} }
    selectedVariableData.value.minItems = undefined
    selectedVariableData.value.maxItems = undefined
  } else {
    delete selectedVariableData.value.items
    delete selectedVariableData.value.minItems
    delete selectedVariableData.value.maxItems
  }
  
  if (newType === 'object') {
    selectedVariableData.value.properties = {}
  } else {
    delete selectedVariableData.value.properties
  }
  
  if (['integer', 'number'].includes(newType)) {
    selectedVariableData.value.min = undefined
    selectedVariableData.value.max = undefined
  } else {
    delete selectedVariableData.value.min
    delete selectedVariableData.value.max
  }
  
  if (newType === 'string') {
    selectedVariableData.value.pattern = undefined
  } else {
    delete selectedVariableData.value.pattern
  }
}

// 确认添加/重命名变量
const confirmAddVariable = () => {
  
  if (!newVariableName.value.trim()) {
    message.warning('请输入变量名')
    return
  }
  
  // 验证变量名格式
  if (!/^[a-zA-Z_][a-zA-Z0-9_]*$/.test(newVariableName.value)) {
    message.error('变量名只能包含字母、数字和下划线，且不能以数字开头')
    return
  }
  
  const variableName = newVariableName.value.trim()
  
  // 处理重命名逻辑
  if (renamingNode.value) {
    const oldPath = renamingNode.value.path
    const pathParts = oldPath.split('.')
    const oldName = pathParts[pathParts.length - 1]
    
    // 如果名称没有改变，直接返回
    if (oldName === variableName) {
      cancelAddVariable()
      return
    }
    
    // 获取变量容器
    const container = getVariableContainer(oldPath)
    if (!container) {
      message.error('无法找到变量容器')
      return
    }
    
    // 检查新名称是否已存在
    if (container[variableName]) {
      message.error('变量名已存在')
      return
    }
    
    // 执行重命名：复制变量数据到新名称，删除旧名称
    const variable = container[oldName]
    // 更新变量的insertText字段为新名称
    if (variable.insertText === getVariableTemplate(oldName)) {
      variable.insertText = getVariableTemplate(variableName)
    }
    container[variableName] = variable
    delete container[oldName]
    
    // 更新选中状态到新路径
    const newPath = pathParts.slice(0, -1).concat([variableName]).join('.')
    selectedKeys.value = [newPath]
    
    // 清除重命名状态
    renamingNode.value = null
    newVariableName.value = ''
    
    // 选择重命名后的变量，不需要强制更新schema
    nextTick(() => {
      onSelectVariable([newPath])
    })
    
    message.success(`变量已重命名为: ${variableName}`)
    return
  }
  
  // 处理新增变量逻辑
  if (!editingNode.value) return
  
  if (editingNode.value.isRoot) {
    try {
      // 检查根级变量名是否重复
      if (varsSchema.value[variableName]) {
        message.error('变量名已存在')
        return
      }
      
      // 创建根级变量 - 使用编辑节点指定的类型
      const variableType = editingNode.value.type || 'string'
      varsSchema.value[variableName] = createDefaultVariable(variableName, variableType)
      
      message.success(`已添加变量: ${variableName}`)
      
      // 先清除编辑状态，再选中新变量
      editingNode.value = null
      newVariableName.value = ''
      
      // 直接选中新创建的变量，不需要强制更新schema
      selectedKeys.value = [variableName]
      nextTick(() => {
        try {
          onSelectVariable([variableName])
        } catch (error) {
          console.error('选择新变量时出错:', error)
        }
      })
    } catch (error) {
      console.error('添加根变量时出错:', error)
      message.error('添加变量失败')
    }
  } else {
    // 处理子变量
    try {
      if (!editingNode.value || !editingNode.value.parentPath) {
        message.error('编辑节点信息不完整')
        return
      }
      
      const parent = getVariableByPath(editingNode.value.parentPath)
      
      if (!parent) {
        message.error('无法找到父级变量')
        return
      }
      
      // 根据父级类型确定子属性的存储位置
      let targetProperties
      if (parent.type === 'object') {
        if (!parent.properties) {
          parent.properties = {}
        }
        targetProperties = parent.properties
      } else if (parent.type === 'object_arr') {
        if (!parent.items) {
          parent.items = { type: 'object', properties: {} }
        }
        if (!parent.items.properties) {
          parent.items.properties = {}
        }
        targetProperties = parent.items.properties
      } else {
        message.error('父级变量类型不支持子变量')
        return
      }
      
      // 检查子变量名是否重复
      if (targetProperties[variableName]) {
        message.error('变量名已存在')
        return
      }
      
      // 创建子变量 - 子变量默认为字符串类型
      // 使用完全独立的对象创建，避免任何引用问题
      const newVariable = {
        type: 'string',
        title: variableName,
        description: '',
        required: false,
        default: '',
        insertText: `{{.${variableName}}}`,
        ui: {
          panel: false,
          order: 10,
          group: '基础信息',
          component: 'input'
        }
      }
      
      // 重建整个 schema 以避免循环引用
      const rawSchema = toRaw(varsSchema.value)
      const newSchema = JSON.parse(JSON.stringify(rawSchema))
      
      // 在新schema中添加变量
      if (parent.type === 'object') {
        if (!newSchema[editingNode.value.parentPath].properties) {
          newSchema[editingNode.value.parentPath].properties = {}
        }
        newSchema[editingNode.value.parentPath].properties[variableName] = newVariable
      } else if (parent.type === 'object_arr') {
        if (!newSchema[editingNode.value.parentPath].items) {
          newSchema[editingNode.value.parentPath].items = { type: 'object', properties: {} }
        }
        if (!newSchema[editingNode.value.parentPath].items.properties) {
          newSchema[editingNode.value.parentPath].items.properties = {}
        }
        newSchema[editingNode.value.parentPath].items.properties[variableName] = newVariable
      }
      
      // 替换整个schema
      varsSchema.value = newSchema
      
      // 选中新创建的子变量
      const newPath = `${editingNode.value.parentPath}.${variableName}`
      
      message.success(`已添加子变量: ${variableName}`)
      
      // 先清除编辑状态，再选中新变量
      editingNode.value = null
      newVariableName.value = ''
      
      // 直接选中新创建的子变量，不需要强制更新schema
      selectedKeys.value = [newPath]
      nextTick(() => {
        try {
          onSelectVariable([newPath])
        } catch (error) {
          console.error('选择新变量时出错:', error)
        }
      })
    } catch (error) {
      console.error('添加子变量时出错:', error)
      message.error('添加子变量失败')
    }
  }
}

// 取消添加变量
const cancelAddVariable = () => {
  if (editingNode.value || renamingNode.value) {
    editingNode.value = null
    renamingNode.value = null
    newVariableName.value = ''
    
    // 不需要强制更新，Vue会自动检测到数据变化
  }
}

// 全局点击处理（自动确认编辑）
const handleGlobalClick = (event) => {
  const inputElement = event.target.closest('.vscode-tree-input')
  const dropdownElement = event.target.closest('.n-dropdown-menu')
  
  // 如果点击的是输入框或下拉菜单，不处理
  if (inputElement || dropdownElement) return
  
  if (editingNode.value || renamingNode.value) {
    confirmAddVariable()
  }
}

// 安全清理循环引用的函数
const cleanCircularReferences = (obj, visited = new WeakSet()) => {
  if (obj === null || typeof obj !== 'object') return obj
  if (visited.has(obj)) return '[Circular Reference]'
  
  visited.add(obj)
  const cleaned = Array.isArray(obj) ? [] : {}
  
  for (const key in obj) {
    if (obj.hasOwnProperty(key)) {
      cleaned[key] = cleanCircularReferences(obj[key], visited)
    }
  }
  
  visited.delete(obj)
  return cleaned
}

// Schema 预览内容 - 显示变量结构定义
const schemaContent = computed(() => {
  try {
    // 先清理循环引用
    const cleanedSchema = cleanCircularReferences(varsSchema.value)
    const schema = { vars_schema: cleanedSchema }
    
    if (previewFormat.value === 'yaml') {
      return YAML.dump(schema, {
        indent: 2,
        lineWidth: -1,
        noRefs: true,
        sortKeys: false
      })
    } else {
      return JSON.stringify(schema, null, 2)
    }
  } catch (error) {
    console.error('生成schema内容时出错:', error)
    return '# Schema 生成失败\n# 错误: ' + error.message
  }
})

// 初始化CodeMirror编辑器
const initSchemaEditor = () => {
  if (!schemaEditorRef.value) return
  
  const extensions = [
    basicSetup,
    EditorView.theme({
      '&': { 
        height: '100%',
        minHeight: '100%'
      },
      '.cm-scroller': { 
        fontFamily: 'Monaco, Menlo, monospace',
        minHeight: '100%'
      },
      '.cm-content': {
        minHeight: '100%',
        padding: '12px'
      },
      '.cm-editor': {
        height: '100%',
        minHeight: '100%'
      },
      '.cm-focused': { outline: 'none' }
    }),
    EditorState.readOnly.of(true),
    EditorView.lineWrapping
  ]
  
  if (previewFormat.value === 'yaml') {
    extensions.push(yaml())
  } else {
    extensions.push(json())
  }
  
  if (schemaEditor) {
    schemaEditor.destroy()
  }
  
  schemaEditor = new EditorView({
    state: EditorState.create({
      doc: schemaContent.value || '{}', // 确保有默认内容
      extensions
    }),
    parent: schemaEditorRef.value
  })
}

// 更新编辑器内容
const updateSchemaEditor = () => {
  try {
    if (!schemaEditor) return
    
    const content = schemaContent.value
    const transaction = schemaEditor.state.update({
      changes: {
        from: 0,
        to: schemaEditor.state.doc.length,
        insert: content
      }
    })
    
    schemaEditor.dispatch(transaction)
  } catch (error) {
    console.error('更新schema编辑器时出错:', error)
  }
}

// 重新初始化编辑器（切换格式时）
const reinitSchemaEditor = async () => {
  await nextTick()
  initSchemaEditor()
}

// 清理损坏的变量数据（移除循环引用）
const cleanVariableData = (obj, visited = new WeakSet()) => {
  if (obj === null || typeof obj !== 'object') return obj
  if (visited.has(obj)) return null // 发现循环引用，返回null
  
  visited.add(obj)
  const cleaned = Array.isArray(obj) ? [] : {}
  
  for (const key in obj) {
    if (obj.hasOwnProperty(key)) {
      const cleanedValue = cleanVariableData(obj[key], visited)
      if (cleanedValue !== null) {
        cleaned[key] = cleanedValue
      }
    }
  }
  
  return cleaned
}

// 加载现有变量定义
const loadVariableDefinitions = async () => {
  try {
    // 从API加载变量定义
    if (props.templateId) {
      const templateId = parseInt(props.templateId)
      const response = await getTemplateExpose({ templateId })
      if (response.data && response.data.data && response.data.data.templateExpose) {
        // 从fieldSchemaJson字段解析JSON
        const fieldSchemaJson = response.data.data.templateExpose.fieldSchemaJson
        if (fieldSchemaJson) {
          varsSchema.value = JSON.parse(fieldSchemaJson)
          console.log('从API加载变量定义:', varsSchema.value)
        } else {
          varsSchema.value = {}
          console.log('使用空的变量定义')
        }
      } else {
        // 如果没有保存的数据，使用空的schema
        varsSchema.value = {}
        console.log('使用空的变量定义')
      }
    } else {
      varsSchema.value = {}
    }
    
    // 重置选择状态
    selectedVariableData.value = null
    selectedKeys.value = []
    editingNode.value = null
    renamingNode.value = null
    newVariableName.value = ''
    
    // 加载或生成测试数据并发送给父组件
    loadTestDataFromStorage()
    emit('test-data-updated', testData.value)
    
  } catch (error) {
    console.error('加载变量定义失败:', error)
    // 如果API加载失败，尝试从localStorage加载作为备用
    try {
      const savedSchema = localStorage.getItem(`template_${props.templateId}_vars_schema`)
      if (savedSchema) {
        const parsed = JSON.parse(savedSchema)
        varsSchema.value = parsed || {}
        console.log('从localStorage加载备用变量定义:', varsSchema.value)
        message.warning('从本地缓存加载变量定义')
      } else {
        varsSchema.value = {}
      }
    } catch (localError) {
      varsSchema.value = {}
      message.error('加载变量定义失败')
    }
    
    // 无论成功还是失败，都要加载并发送测试数据
    loadTestDataFromStorage()
    emit('test-data-updated', testData.value)
  }
}

// 处理保存点击
const handleSaveClick = () => {
  console.log('保存按钮被点击', { templateId: props.templateId, saving: saving.value })
  if (!props.templateId) {
    message.error('模板ID不存在，无法保存')
    return
  }
  if (saving.value) {
    console.log('正在保存中，跳过重复操作')
    return
  }
  saveVariables()
}

// 保存变量定义
const saveVariables = async () => {
  console.log('保存变量定义被调用', { templateId: props.templateId, varsSchema: varsSchema.value })
  saving.value = true
  try {
    // 同步当前编辑的变量数据到schema
    if (selectedVariableData.value && selectedKeys.value.length > 0) {
      const path = selectedKeys.value[0]
      updateVariableInSchema(path, selectedVariableData.value)
    }
    
    // 清理循环引用
    const cleanedSchema = cleanCircularReferences(varsSchema.value)
    
    // 调用API保存到后端
    const templateId = parseInt(props.templateId)
    await setTemplateExpose({
      templateId,
      varsSchema: cleanedSchema,
      version: '1.0'
    })
    
    // 同时保存到本地存储作为备份
    localStorage.setItem(`template_${props.templateId}_vars_schema`, JSON.stringify(cleanedSchema))
    
    message.success('变量定义保存成功')
    
    // 使用已保存的测试数据（不自动重新生成，让用户控制测试数据）
    loadTestDataFromStorage()
    emit('test-data-updated', testData.value)
    
    // 通知父组件刷新用户变量
    emit('variables-saved')
  } catch (error) {
    console.error('保存失败:', error)
    message.error('保存失败：' + (error.message || '未知错误'))
  } finally {
    saving.value = false
  }
}

// 更新变量到schema
const updateVariableInSchema = (path, variableData) => {
  const variable = getVariableByPath(path)
  
  if (variable && variableData) {
    // 只更新基本字段，避免覆盖properties导致循环引用
    const { path: _, properties, ...configData } = variableData
    
    // 安全地更新基本属性
    Object.keys(configData).forEach(key => {
      if (key !== 'properties') {
        variable[key] = configData[key]
      }
    })
    
    // properties字段单独处理，避免循环引用
    // 如果原variable没有properties但新数据有，才创建空的properties
    if (!variable.properties && variableData.type === 'object') {
      variable.properties = {}
    }
  }
}

// 监听变量数据变化，实时同步到schema
watch(selectedVariableData, (newData) => {
  if (newData && selectedKeys.value.length > 0) {
    console.log('变量数据更新:', selectedKeys.value[0], newData.type)
    updateVariableInSchema(selectedKeys.value[0], newData)
    
    // 注意：不需要强制更新varsSchema，因为updateVariableInSchema已经直接修改了对象
    // 过度的强制更新可能导致表单重复渲染
  }
}, { deep: true })

// 拖拽调整高度
const startResize = (e) => {
  isResizing.value = true
  const startY = e.clientY
  const startHeight = drawerHeight.value

  const onMouseMove = (e) => {
    if (!isResizing.value) return
    
    const deltaY = startY - e.clientY // 向上为正，向下为负
    const newHeight = Math.min(Math.max(startHeight + deltaY, minHeight), maxHeight)
    drawerHeight.value = newHeight
  }

  const onMouseUp = () => {
    isResizing.value = false
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
    document.body.style.userSelect = ''
    document.body.style.cursor = ''
  }

  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'ns-resize'
}

// 初始化面板宽度
const initPanelWidths = () => {
  const layoutElement = document.querySelector('.expose-layout')
  if (layoutElement) {
    totalLayoutWidth.value = layoutElement.offsetWidth
    rightPanelWidth.value = Math.floor(totalLayoutWidth.value / 3)
    leftPanelWidth.value = 250
  }
}

// 左侧面板拖拽调整
const startLeftResize = (e) => {
  isLeftResizing.value = true
  const startX = e.clientX
  const startWidth = leftPanelWidth.value
  
  const onMouseMove = (e) => {
    if (!isLeftResizing.value) return
    
    const deltaX = e.clientX - startX
    const newWidth = Math.min(Math.max(startWidth + deltaX, 100), totalLayoutWidth.value - 400) // 最小100px，保证其他面板至少400px
    leftPanelWidth.value = newWidth
  }
  
  const onMouseUp = () => {
    isLeftResizing.value = false
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
    document.body.style.userSelect = ''
    document.body.style.cursor = ''
  }
  
  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'ew-resize'
}

// 右侧面板拖拽调整
const startRightResize = (e) => {
  isRightResizing.value = true
  const startX = e.clientX
  const startWidth = rightPanelWidth.value
  
  const onMouseMove = (e) => {
    if (!isRightResizing.value) return
    
    const deltaX = startX - e.clientX // 右侧面板向左拖拽为正
    const newWidth = Math.min(Math.max(startWidth + deltaX, 150), totalLayoutWidth.value - leftPanelWidth.value - 200) // 最小150px，保证左侧和中间至少200px
    rightPanelWidth.value = newWidth
  }
  
  const onMouseUp = () => {
    isRightResizing.value = false
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
    document.body.style.userSelect = ''
    document.body.style.cursor = ''
  }
  
  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'ew-resize'
}

// 紧急清理函数
const emergencyCleanup = () => {
  try {
    // 清空所有可能的本地存储
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key && key.includes('vars_schema')) {
        localStorage.removeItem(key)
      }
    }
    
    // 重置所有响应式数据
    varsSchema.value = {}
    selectedVariableData.value = null
    selectedKeys.value = []
    editingNode.value = null
    renamingNode.value = null
    newVariableName.value = ''
    
    // 强制触发重新渲染
    nextTick(() => {
      console.log('紧急清理完成，数据已重置')
    })
    
    message.success('所有变量数据已重置，可以重新开始创建变量')
  } catch (error) {
    console.error('紧急清理失败:', error)
    message.error('重置失败')
  }
}

// 监听抽屉显示状态
watch(visible, async (show) => {
  if (show) {
    console.log('变量定义抽屉打开', { templateId: props.templateId, type: typeof props.templateId })
    // 加载变量定义数据
    loadVariableDefinitions()
    await nextTick()
    initPanelWidths()
    initSchemaEditor()
  }
})

// 监听预览格式变化
watch(previewFormat, () => {
  reinitSchemaEditor()
})

// Schema预览格式切换监听

// 监听varsSchema变化 - 恢复深度监听，现在循环引用问题已解决
watch(varsSchema, (newSchema) => {
  try {
    console.log('varsSchema变化，更新预览')
    updateSchemaEditor()
  } catch (error) {
    console.error('varsSchema watcher 执行时出错:', error)
    // 如果遇到错误，尝试清理循环引用后重试
    try {
      updateSchemaEditor()
    } catch (retryError) {
      console.error('重试更新预览失败:', retryError)
    }
  }
}, { deep: true })

// 测试数据生成和管理函数
const generateTestData = (schema) => {
  const data = {}
  
  if (!schema || typeof schema !== 'object') {
    return data
  }
  
  Object.entries(schema).forEach(([key, variable]) => {
    if (!variable || typeof variable !== 'object') return
    
    switch (variable.type) {
      case 'string':
        data[key] = variable.default || `示例${key}`
        break
      case 'integer':
        data[key] = variable.default || 42
        break
      case 'number':
        data[key] = variable.default || 3.14
        break
      case 'boolean':
        data[key] = variable.default !== undefined ? variable.default : true
        break
      case 'array':
        data[key] = variable.default || ['item1', 'item2']
        break
      case 'object':
        if (variable.properties) {
          data[key] = generateTestData(variable.properties)
        } else {
          data[key] = {}
        }
        break
      case 'object_arr':
        if (variable.items && variable.items.properties) {
          const itemData = generateTestData(variable.items.properties)
          data[key] = [itemData, { ...itemData }] // 生成两个示例项
        } else {
          data[key] = []
        }
        break
      case 'enum':
        data[key] = variable.enum && variable.enum.length > 0 ? variable.enum[0] : variable.default || ''
        break
      case 'secret':
        data[key] = variable.default || '***保密信息***'
        break
      default:
        data[key] = variable.default || ''
    }
  })
  
  return data
}

const showTestDataModal = () => {
  // 先从localStorage加载已保存的测试数据，如果没有则生成新的
  loadTestDataFromStorage()
  showTestDataModalFlag.value = true
  
  // 下次tick时初始化编辑器
  nextTick(() => {
    initTestDataEditor()
  })
}

// 从localStorage加载测试数据
const loadTestDataFromStorage = () => {
  const testDataKey = `template_studio_${props.templateId}_testdata`
  const savedTestData = localStorage.getItem(testDataKey)
  
  if (savedTestData) {
    try {
      testData.value = JSON.parse(savedTestData)
      console.log('从localStorage加载测试数据:', testData.value)
    } catch (e) {
      console.error('解析保存的测试数据失败:', e)
      // 解析失败则重新生成
      testData.value = generateTestData(varsSchema.value)
    }
  } else {
    // 没有保存的数据，生成新的测试数据
    testData.value = generateTestData(varsSchema.value)
  }
}

// 保存测试数据到localStorage
const saveTestDataToStorage = () => {
  try {
    // 从编辑器获取当前数据
    if (testDataEditor.value?.codemirror) {
      const editorContent = testDataEditor.value.codemirror.state.doc.toString()
      const parsedData = JSON.parse(editorContent)
      testData.value = parsedData
    }
    
    const testDataKey = `template_studio_${props.templateId}_testdata`
    localStorage.setItem(testDataKey, JSON.stringify(testData.value))
    
    // 同时更新模板预览的数据
    emit('test-data-updated', testData.value)
    
    message.success('测试数据已保存')
  } catch (error) {
    console.error('保存测试数据失败:', error)
    message.error('保存失败: ' + error.message)
  }
}

const initTestDataEditor = () => {
  if (!testDataEditor.value) return
  
  // 清除现有编辑器
  if (testDataEditor.value.codemirror) {
    testDataEditor.value.codemirror.destroy()
  }
  
  const state = EditorState.create({
    doc: JSON.stringify(testData.value, null, 2),
    extensions: [
      basicSetup,
      json(),
      EditorView.theme({
        '&': {
          fontSize: '14px',
          fontFamily: 'Monaco, Menlo, "Ubuntu Mono", Consolas, monospace'
        },
        '.cm-content': {
          padding: '12px',
          minHeight: '200px'
        },
        '.cm-editor': {
          borderRadius: '6px',
          border: '1px solid #e0e0e6'
        },
        '.cm-focused': {
          outline: 'none',
          borderColor: '#18a058'
        }
      }),
      EditorView.updateListener.of((update) => {
        if (update.docChanged) {
          try {
            const newData = JSON.parse(update.state.doc.toString())
            testData.value = newData
          } catch (error) {
            // JSON格式错误，暂不处理
          }
        }
      })
    ]
  })
  
  const view = new EditorView({
    state,
    parent: testDataEditor.value
  })
  
  testDataEditor.value.codemirror = view
}

const copyTestData = () => {
  const dataStr = JSON.stringify(testData.value, null, 2)
  navigator.clipboard.writeText(dataStr).then(() => {
    message.success('测试数据已复制到剪贴板')
  }).catch(() => {
    message.error('复制失败')
  })
}

const regenerateTestData = () => {
  testData.value = generateTestData(varsSchema.value)
  if (testDataEditor.value?.codemirror) {
    const newDoc = JSON.stringify(testData.value, null, 2)
    testDataEditor.value.codemirror.dispatch({
      changes: {
        from: 0,
        to: testDataEditor.value.codemirror.state.doc.length,
        insert: newDoc
      }
    })
  }
  message.success('测试数据已重新生成，请点击"保存数据"以应用更改')
}

// 变量分析相关函数
const showVariableAnalysisModal = async () => {
  if (!props.templateId) {
    message.error('模板ID不能为空')
    return
  }
  
  analyzingVariables.value = true
  try {
    const result = await analyzeTemplateVariables(props.templateId)
    console.log('分析结果原始数据:', result)
    console.log('分析结果data部分:', result.data)
    // 修正：应该使用result.data.data，因为后端返回格式是{code, message, data: {实际数据}}
    analysisResult.value = result.data.data
    console.log('修正后的analysisResult:', analysisResult.value)
    console.log('缺失变量数量:', analysisResult.value?.missingVariables?.length)
    console.log('检测变量数量:', analysisResult.value?.detectedVariables?.length)
    showVariableAnalysisModalFlag.value = true
  } catch (error) {
    console.error('分析变量失败:', error)
    message.error('分析变量失败: ' + (error.message || '未知错误'))
  } finally {
    analyzingVariables.value = false
  }
}

// 添加单个检测到的变量（智能类型推断）
const addDetectedVariable = (detectedVar) => {
  if (!detectedVar) return
  
  const currentSchema = toRaw(varsSchema.value)
  const newSchema = JSON.parse(JSON.stringify(currentSchema))
  
  if (newSchema[detectedVar.name]) {
    message.warning(`变量 "${detectedVar.name}" 已存在`)
    return
  }
  
  // 智能类型推断
  const smartType = getSuggestedType(detectedVar)
  
  // 创建新变量
  newSchema[detectedVar.name] = {
    type: smartType,
    title: detectedVar.name,
    description: `模板变量: 在${detectedVar.files.length}个文件中使用`,
    required: false,
    default: getDefaultValueForType(smartType),
    insertText: `{{.${detectedVar.name}}}`,
    ui: {
      panel: false,
      order: 10,
      group: '模板变量',
      component: getDefaultComponent(smartType)
    }
  }
  
  // 为特定类型添加额外字段
  if (smartType === 'object') {
    newSchema[detectedVar.name].properties = {}
  } else if (smartType === 'array') {
    newSchema[detectedVar.name].items = { type: 'string' }
  } else if (smartType === 'object_arr') {
    newSchema[detectedVar.name].items = {
      type: 'object',
      properties: {}
    }
  }
  
  varsSchema.value = newSchema
  message.success(`已添加变量 "${detectedVar.name}" (${smartType}类型)`)
  // 自动展开新添加的变量
  expandedKeys.value = [...expandedKeys.value, detectedVar.name]
  
  // 实时更新缺失变量列表
  if (analysisResult.value?.missingVariables) {
    analysisResult.value.missingVariables = analysisResult.value.missingVariables.filter(v => v.name !== detectedVar.name)
  }
}

// 添加所有检测到的变量
const addAllDetectedVariables = (variablesToAdd) => {
  if (!variablesToAdd || variablesToAdd.length === 0) {
    message.warning('没有可添加的变量')
    return
  }
  
  const currentSchema = toRaw(varsSchema.value)
  const newSchema = JSON.parse(JSON.stringify(currentSchema))
  
  let addedCount = 0
  const addedNames = []
  
  variablesToAdd.forEach(detectedVar => {
    if (!newSchema[detectedVar.name]) {
      // 智能类型推断
      const smartType = getSuggestedType(detectedVar)
      
      // 创建新变量
      newSchema[detectedVar.name] = {
        type: smartType,
        title: detectedVar.name,
        description: `模板变量: 在${detectedVar.files.length}个文件中使用`,
        required: false,
        default: getDefaultValueForType(smartType),
        insertText: `{{.${detectedVar.name}}}`,
        ui: {
          panel: false,
          order: 10,
          group: '模板变量',
          component: getDefaultComponent(smartType)
        }
      }
      
      // 为特定类型添加额外字段
      if (smartType === 'object') {
        newSchema[detectedVar.name].properties = {}
      } else if (smartType === 'array') {
        newSchema[detectedVar.name].items = { type: 'string' }
      } else if (smartType === 'object_arr') {
        newSchema[detectedVar.name].items = {
          type: 'object',
          properties: {}
        }
      }
      
      addedCount++
      addedNames.push(detectedVar.name)
    }
  })
  
  if (addedCount > 0) {
    varsSchema.value = newSchema
    message.success(`已添加 ${addedCount} 个变量到变量树`)
    // 自动展开新添加的变量
    expandedKeys.value = [...expandedKeys.value, ...addedNames]
    
    // 实时更新缺失变量列表
    if (analysisResult.value?.missingVariables) {
      analysisResult.value.missingVariables = analysisResult.value.missingVariables.filter(v => !addedNames.includes(v.name))
    }
  } else {
    message.info('所有检测到的变量都已存在')
  }
}


// 根据上下文和建议获取推荐类型
const getSuggestedType = (detectedVar) => {
  if (!detectedVar.contexts || detectedVar.contexts.length === 0) {
    return detectedVar.type === 'unknown' ? 'string' : detectedVar.type
  }
  
  // 分析上下文模式
  const contexts = detectedVar.contexts
  const hasRangeContext = contexts.some(c => c.includes('range'))
  const hasIfContext = contexts.some(c => c.includes('if'))
  const hasIndexContext = contexts.some(c => c.includes('['))
  const hasNestedContext = contexts.some(c => c.includes('.') && !c.match(/\{\{\.[^.]+\}\}$/))
  
  if (hasRangeContext) {
    return 'array' // range 通常用于数组遍历
  } else if (hasIfContext) {
    return 'boolean' // if 通常用于布尔判断
  } else if (hasIndexContext) {
    return 'array' // 索引访问通常是数组
  } else if (hasNestedContext) {
    return 'object' // 嵌套访问通常是对象
  }
  
  return detectedVar.type === 'unknown' ? 'string' : detectedVar.type
}

// 删除单个未使用变量
const deleteUnusedVariable = (varName) => {
  if (!varName) return
  
  const currentSchema = toRaw(varsSchema.value)
  const newSchema = JSON.parse(JSON.stringify(currentSchema))
  
  if (newSchema[varName]) {
    delete newSchema[varName]
    varsSchema.value = newSchema
    message.success(`已删除未使用变量 "${varName}"`)
    
    // 清理选择状态
    if (selectedKeys.value.includes(varName)) {
      selectedKeys.value = []
      selectedVariableData.value = null
    }
    
    // 从展开的键中移除
    expandedKeys.value = expandedKeys.value.filter(key => key !== varName)
    
    // 实时更新未使用变量列表
    if (analysisResult.value?.unusedVariables) {
      analysisResult.value.unusedVariables = analysisResult.value.unusedVariables.filter(name => name !== varName)
    }
  } else {
    message.warning(`变量 "${varName}" 不存在`)
  }
}

// 删除所有未使用变量
const deleteAllUnusedVariables = (unusedVars) => {
  if (!unusedVars || unusedVars.length === 0) {
    message.warning('没有未使用的变量需要删除')
    return
  }
  
  const currentSchema = toRaw(varsSchema.value)
  const newSchema = JSON.parse(JSON.stringify(currentSchema))
  let deletedCount = 0
  const deletedNames = []
  
  unusedVars.forEach(varName => {
    if (newSchema[varName]) {
      delete newSchema[varName]
      deletedCount++
      deletedNames.push(varName)
    }
  })
  
  if (deletedCount > 0) {
    varsSchema.value = newSchema
    message.success(`已删除 ${deletedCount} 个未使用的变量`)
    
    // 清理选择状态
    if (selectedKeys.value.some(key => deletedNames.includes(key))) {
      selectedKeys.value = []
      selectedVariableData.value = null
    }
    
    // 从展开的键中移除
    expandedKeys.value = expandedKeys.value.filter(key => !deletedNames.includes(key))
    
    // 实时更新未使用变量列表
    if (analysisResult.value?.unusedVariables) {
      analysisResult.value.unusedVariables = analysisResult.value.unusedVariables.filter(name => !deletedNames.includes(name))
    }
  } else {
    message.info('没有可删除的变量')
  }
}

// 清理未使用的变量（保留原函数名以兼容统计页面）
const cleanupUnusedVariables = () => {
  deleteAllUnusedVariables(analysisResult.value?.unusedVariables)
}

// 优化所有变量类型
const optimizeAllVariableTypes = () => {
  if (!analysisResult.value?.detectedVariables || analysisResult.value.detectedVariables.length === 0) {
    message.warning('没有检测到的变量需要优化')
    return
  }
  
  const currentSchema = toRaw(varsSchema.value)
  const newSchema = JSON.parse(JSON.stringify(currentSchema))
  let optimizedCount = 0
  
  analysisResult.value.detectedVariables.forEach(detectedVar => {
    if (newSchema[detectedVar.name]) {
      const suggestedType = getSuggestedType(detectedVar)
      if (newSchema[detectedVar.name].type !== suggestedType) {
        newSchema[detectedVar.name].type = suggestedType
        newSchema[detectedVar.name].description = `自动优化: ${detectedVar.suggestions}`
        newSchema[detectedVar.name].default = getDefaultValueForType(suggestedType)
        newSchema[detectedVar.name].ui.component = getDefaultComponent(suggestedType)
        
        // 为特定类型添加额外字段
        if (suggestedType === 'object' && !newSchema[detectedVar.name].properties) {
          newSchema[detectedVar.name].properties = {}
        } else if (suggestedType === 'array' && !newSchema[detectedVar.name].items) {
          newSchema[detectedVar.name].items = { type: 'string' }
        } else if (suggestedType === 'object_arr' && !newSchema[detectedVar.name].items) {
          newSchema[detectedVar.name].items = {
            type: 'object',
            properties: {}
          }
        }
        
        optimizedCount++
      }
    }
  })
  
  if (optimizedCount > 0) {
    varsSchema.value = newSchema
    message.success(`已优化 ${optimizedCount} 个变量的类型`)
  } else {
    message.info('所有变量类型都已是最优')
  }
}

const getDefaultValueForType = (type) => {
  const typeMap = {
    string: '',
    integer: 0,
    number: 0.0,
    boolean: false,
    array: [],
    object: {},
    object_arr: [],
    unknown: ''
  }
  return typeMap[type] || ''
}

// 预览面板操作函数
// 1. 导入文件
const importFromFile = () => {
  fileInputRef.value?.click()
}

// 2. 处理文件导入
const handleFileImport = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  try {
    const text = await file.text()
    let importedData = null
    
    // 根据文件扩展名或内容判断格式
    if (file.name.endsWith('.json') || text.trim().startsWith('{')) {
      importedData = JSON.parse(text)
    } else if (file.name.endsWith('.yaml') || file.name.endsWith('.yml')) {
      importedData = YAML.load(text)
    } else {
      // 尝试两种格式
      try {
        importedData = JSON.parse(text)
      } catch {
        importedData = YAML.load(text)
      }
    }
    
    // 验证数据结构
    if (importedData && typeof importedData === 'object') {
      // 如果有vars_schema字段，使用它；否则直接使用整个对象
      const schemaData = importedData.vars_schema || importedData
      
      // 保存待导入的数据并显示确认对话框
      pendingImportData.value = schemaData
      showImportConfirmModal.value = true
    } else {
      message.error('文件格式不正确，请确保是有效的JSON或YAML格式')
    }
  } catch (error) {
    console.error('导入文件时出错:', error)
    message.error('文件解析失败，请检查文件格式')
  }
  
  // 清空input值以便重新选择同一文件
  event.target.value = ''
}

// 3. 显示导出对话框
const showExportDialog = () => {
  showExportModalFlag.value = true
}

// 4. 导出文件
const exportFile = (format) => {
  try {
    const schema = { vars_schema: varsSchema.value }
    let content = ''
    let filename = ''
    let mimeType = ''
    
    if (format === 'json') {
      content = JSON.stringify(schema, null, 2)
      filename = `template_vars_schema_${Date.now()}.json`
      mimeType = 'application/json'
    } else if (format === 'yaml') {
      content = YAML.dump(schema, {
        indent: 2,
        lineWidth: -1,
        noRefs: true,
        sortKeys: false
      })
      filename = `template_vars_schema_${Date.now()}.yaml`
      mimeType = 'application/x-yaml'
    }
    
    // 创建下载链接
    const blob = new Blob([content], { type: mimeType })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    
    showExportModalFlag.value = false
    message.success(`已导出为 ${format.toUpperCase()} 格式`)
  } catch (error) {
    console.error('导出文件时出错:', error)
    message.error('导出失败')
  }
}

// 5. 复制到剪贴板
const copyToClipboard = async () => {
  try {
    const content = schemaContent.value
    await navigator.clipboard.writeText(content)
    message.success('已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    message.error('复制失败')
  }
}

// 6. 格式化内容
const formatContent = () => {
  try {
    updateSchemaEditor()
    message.success('内容已格式化')
  } catch (error) {
    console.error('格式化失败:', error)
    message.error('格式化失败')
  }
}

// 7. 同步到树
const syncToTree = () => {
  try {
    // 获取编辑器内容
    if (!schemaEditor) {
      message.warning('预览编辑器未初始化')
      return
    }
    
    const editorContent = schemaEditor.state.doc.toString()
    let parsedData = null
    
    // 根据当前格式解析
    if (previewFormat.value === 'json') {
      parsedData = JSON.parse(editorContent)
    } else {
      parsedData = YAML.load(editorContent)
    }
    
    // 验证并应用
    if (parsedData && typeof parsedData === 'object') {
      const schemaData = parsedData.vars_schema || parsedData
      
      // 保存待同步的数据并显示确认对话框
      pendingSyncData.value = schemaData
      showSyncConfirmModal.value = true
    } else {
      message.error('预览内容格式不正确')
    }
  } catch (error) {
    console.error('同步失败:', error)
    message.error('内容解析失败，请检查格式')
  }
}

// 确认导入
const confirmImport = async () => {
  if (pendingImportData.value) {
    varsSchema.value = pendingImportData.value
    selectedVariableData.value = null
    selectedKeys.value = []
    showImportConfirmModal.value = false
    pendingImportData.value = null
    
    // 导入后自动保存
    message.success('Schema导入成功，正在保存...')
    await saveVariables()
    // saveVariables 中已经会触发 emit('variables-saved')
  }
}

// 确认同步
const confirmSync = async () => {
  if (pendingSyncData.value) {
    varsSchema.value = pendingSyncData.value
    selectedVariableData.value = null
    selectedKeys.value = []
    showSyncConfirmModal.value = false
    pendingSyncData.value = null
    
    // 同步后自动保存
    message.success('已同步到变量树，正在保存...')
    await saveVariables()
    // saveVariables 中已经会触发 emit('variables-saved')
  }
}

// 组件挂载时添加全局事件监听
onMounted(() => {
  document.addEventListener('click', handleGlobalClick, true)
})

// 组件卸载时清理
onUnmounted(() => {
  if (schemaEditor) {
    schemaEditor.destroy()
    schemaEditor = null
  }
  document.removeEventListener('click', handleGlobalClick, true)
})

// 暴露方法给父组件
defineExpose({
  showTestDataModal
})
</script>

<style scoped>
.resize-handle {
  position: relative;
  height: 12px;
  cursor: ns-resize;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 8px;
}

.resize-handle:hover {
  background-color: rgba(24, 144, 255, 0.1);
}

.handle-bar {
  width: 60px;
  height: 4px;
  background: #d9d9d9;
  border-radius: 2px;
  transition: background-color 0.2s;
}

.resize-handle:hover .handle-bar {
  background: #1890ff;
}

.expose-layout {
  display: flex;
  height: calc(100% - 20px);
  gap: 0;
  position: relative;
}

.left-panel,
.center-panel,
.right-panel {
  background: #f8f8f8;
  border-radius: 6px;
  padding: 16px;
  overflow-y: auto;
  flex-shrink: 0;
}

/* 预览面板头部样式 */
.panel-header {
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #e0e0e0;
}

.panel-title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.panel-header .panel-title {
  font-weight: 600;
  color: #333;
  margin: 0;
  font-size: 14px;
}

.panel-actions {
  display: flex;
  justify-content: flex-end;
  flex-wrap: wrap;
}

.panel-actions .n-space {
  flex-wrap: wrap;
}

.left-panel {
  min-width: 100px;
  margin-right: 4px;
}

.center-panel {
  min-width: 200px;
  margin: 0 4px;
  flex: 1;
}

.right-panel {
  min-width: 150px;
  margin-left: 4px;
}

.panel-resizer {
  width: 8px;
  background: transparent;
  cursor: ew-resize;
  position: relative;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
}

.panel-resizer:hover {
  background-color: rgba(24, 144, 255, 0.1);
}

.panel-resizer::before {
  content: '';
  width: 2px;
  height: 20px;
  background: #d9d9d9;
  border-radius: 1px;
  transition: background-color 0.2s;
}

.panel-resizer:hover::before {
  background: #1890ff;
}

.left-resizer {
  order: 1;
}

.center-panel {
  order: 2;
}

.right-resizer {
  order: 3;
}

.right-panel {
  order: 4;
}

.left-panel {
  order: 0;
}

.panel-title {
  font-weight: 600;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.file-list {
  margin-bottom: 16px;
}

.file-item {
  display: flex;
  align-items: center;
  padding: 8px;
  margin-bottom: 4px;
  background: white;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.file-item:hover {
  background: #e6f7ff;
}

.file-item.active {
  background: #1890ff;
  color: white;
}

.file-icon {
  margin-right: 8px;
}

.usage-section {
  border-top: 1px solid #e0e0e0;
  padding-top: 12px;
}

.section-title {
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #666;
}

.usage-list {
  max-height: 150px;
  overflow-y: auto;
}

.usage-item {
  display: flex;
  justify-content: space-between;
  padding: 4px 8px;
  margin-bottom: 2px;
  background: white;
  border-radius: 3px;
  font-size: 12px;
}

.var-name {
  color: #1890ff;
  font-family: monospace;
}

.var-count {
  color: #666;
}

.variables-list {
  margin-bottom: 16px;
  max-height: 200px;
  overflow-y: auto;
}

.variable-card {
  background: white;
  border-radius: 4px;
  padding: 12px;
  margin-bottom: 8px;
  cursor: pointer;
  border: 1px solid #e0e0e0;
  transition: all 0.2s;
}

.variable-card:hover {
  border-color: #1890ff;
}

.variable-card.active {
  border-color: #1890ff;
  background: #f0f9ff;
}

.variable-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.variable-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.variable-name {
  font-weight: 600;
}

.variable-type {
  font-size: 12px;
  padding: 2px 6px;
  background: #f0f0f0;
  border-radius: 3px;
  color: #666;
}

.variable-desc {
  font-size: 12px;
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.variable-form {
  border-top: 1px solid #e0e0e0;
  padding-top: 16px;
}

.preview-content {
  height: calc(100% - 50px);
  overflow: hidden;
}

.code-preview {
  height: 100%;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  overflow: hidden;
}

.code-preview :deep(.cm-editor) {
  height: 100% !important;
  min-height: 100%;
}

.code-preview :deep(.cm-scroller) {
  font-size: 12px;
  min-height: 100%;
}

.code-preview :deep(.cm-content) {
  min-height: 100%;
  padding: 12px;
}

.code-preview :deep(.cm-focused) {
  outline: none;
}

.json-preview {
  font-family: monospace;
  font-size: 12px;
  background: white;
  padding: 12px;
  border-radius: 4px;
  border: 1px solid #e0e0e0;
  margin: 0;
  white-space: pre-wrap;
}


/* 变量树样式 */
.variable-tree {
  height: calc(100% - 50px);
  overflow-y: auto;
}

.var-icon {
  margin-right: 4px;
}

.var-icon.var-string {
  color: #1890ff; /* 蓝色 - 文本 */
}

.var-icon.var-integer,
.var-icon.var-number {
  color: #52c41a; /* 绿色 - 数字 */
}

.var-icon.var-boolean {
  color: #fa8c16; /* 橙色 - 布尔值 */
}

.var-icon.var-array {
  color: #f5222d; /* 红色 - 数组 */
}

.var-icon.var-object {
  color: #722ed1; /* 紫色 - 对象 */
}

.var-icon.var-enum {
  color: #13c2c2; /* 青色 - 枚举 */
}

.var-icon.var-secret {
  color: #eb2f96; /* 粉色 - 密码 */
}

.node-actions {
  opacity: 0;
  transition: opacity 0.2s;
}

.n-tree-node:hover .node-actions {
  opacity: 1;
}

.tree-node-content {
  display: flex;
  align-items: center;
  gap: 6px;
  width: 100%;
}

.node-title {
  flex: 1;
}

/* 表单样式 */
.form-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px;
  background: #fafafa;
  border-radius: 6px;
}

.var-path {
  font-family: monospace;
  color: #666;
  font-size: 12px;
}

.form-section-title {
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
  padding-bottom: 4px;
  border-bottom: 1px solid #e8e8e8;
}

.empty-selection {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.empty-tree-hint {
  padding: 32px;
  color: #888;
  text-align: center;
  user-select: none;
  cursor: context-menu;
  font-size: 14px;
  border: 2px dashed #e0e0e0;
  border-radius: 6px;
  margin: 16px;
  transition: all 0.2s;
}

.empty-tree-hint:hover {
  border-color: #1890ff;
  color: #1890ff;
  background: rgba(24, 144, 255, 0.05);
}

/* VSCode风格的树输入框样式 */
:deep(.vscode-tree-input) {
  width: 100%;
  height: 22px;
  padding: 1px 4px;
  font-size: 13px;
  font-family: 'Segoe UI', 'Consolas', 'Monaco', monospace;
  background: #ffffff;
  border: 1px solid #007acc;
  border-radius: 0;
  outline: none;
  color: #333333;
  line-height: 18px;
  box-shadow: 0 0 0 1px #007acc;
  margin: 0;
  display: block;
  box-sizing: border-box;
}

:deep(.vscode-tree-input:focus) {
  border-color: #007acc;
  box-shadow: 0 0 0 2px rgba(0, 122, 204, 0.25);
}

:deep(.vscode-tree-input::placeholder) {
  color: #999999;
  font-style: italic;
}

/* 对象数组编辑器样式 */
.object-structure-hint {
  margin-bottom: 8px;
}

.object-properties-editor {
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  padding: 12px;
  background: #fafafa;
  min-height: 120px;
}

.empty-properties-hint {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100px;
  color: #999;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.empty-properties-hint:hover {
  background: #f0f0f0;
  color: #666;
}

.properties-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.property-item {
  padding: 8px 12px;
  background: white;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  transition: all 0.2s;
}

.property-item:hover {
  border-color: #d4edda;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.property-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.property-name {
  font-weight: 600;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  color: #333;
  flex: 1;
}

.property-description {
  font-size: 12px;
  color: #666;
  padding-left: 4px;
}
</style>