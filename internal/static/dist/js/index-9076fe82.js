import{d as defineComponent,I as createElementBlock,a0 as createBaseVNode,R as normalizeClass,ab as createTextVNode,a1 as toDisplayString,J as renderSlot,z as openBlock,u as unref,r as ref,c as computed,F as Fragment,M as renderList,a5 as withModifiers,E as resolveComponent,f as createVNode,L as withCtx,w as watch,K as createCommentVNode,ac as createStaticVNode,n as nextTick,i as isRef,Q as normalizeStyle,H as withDirectives,$ as vShow,A as createBlock,b as onMounted,o as onUnmounted,ad as useSlots,a6 as Transition}from"./vue-vendor-31320f47.js";import{_ as _export_sfc,j as useTexts,aA as API_BASE_URL,q as useToast,J as getCurrentLocale,v as logger}from"./index-38202212.js";import{H as HighlightJS}from"./highlight-core-d3d00c28.js";import{b as getImageFormatOptions}from"./filterOptions-1ed3ac9a.js";import{g as getDocsNavItems,C as CODE_TAB_OPTIONS}from"./navigationOptions-d958d3cc.js";import"./editor-markdown-0b2e1ce5.js";import"./axios-c2a59b9f.js";import"./spark-md5-d0a97131.js";import"./highlight-common-d6c232ae.js";import"./comlink-c7be92b7.js";/* empty css                 */import"./vueuse-2480d07f.js";import"./marked-e9418aa8.js";import"./dompurify-76a35038.js";import"./vue-table-79795ca5.js";import"./highlight-other-84e6ee74.js";const _sfc_main$j=defineComponent({name:"ApiSection",props:{id:{type:String,required:!0},title:{type:String,required:!0},iconClass:{type:String,required:!0}}});const _hoisted_1$j=["id"],_hoisted_2$i={class:"mb-4 flex items-center text-xl font-bold"};function _sfc_render$5(_ctx,_cache,$props,$setup,$data,$options){return openBlock(),createElementBlock("div",{id:_ctx.id,class:"cyber-panel api-section mb-6 rounded-lg p-6"},[createBaseVNode("h3",_hoisted_2$i,[createBaseVNode("i",{class:normalizeClass(["mr-2",[_ctx.iconClass]])},null,2),createTextVNode(toDisplayString(_ctx.title),1)]),renderSlot(_ctx.$slots,"default",{},void 0,!0)],8,_hoisted_1$j)}const ApiSection=_export_sfc(_sfc_main$j,[["render",_sfc_render$5],["__scopeId","data-v-292d2173"]]),_hoisted_1$i={class:"cyber-tester-header"},_hoisted_2$h={class:"header-left"},_hoisted_3$g={class:"endpoint-display"},_hoisted_4$e={class:"endpoint-url"},_hoisted_5$c={class:"header-right"},_hoisted_6$c=["disabled"],_hoisted_7$c={class:"btn-content"},_sfc_main$i=defineComponent({name:"TesterHeader",__name:"TesterHeader",props:{currentDomain:{},canSend:{type:Boolean},isLoading:{type:Boolean}},emits:["send-request"],setup(__props){const{$t:$t2}=useTexts();return(_ctx,_cache)=>(openBlock(),createElementBlock("div",_hoisted_1$i,[createBaseVNode("div",_hoisted_2$h,[createBaseVNode("div",_hoisted_3$g,[_cache[1]||(_cache[1]=createBaseVNode("div",{class:"method-badge"},"POST",-1)),createBaseVNode("div",_hoisted_4$e,[createBaseVNode("code",null,toDisplayString(__props.currentDomain)+"/api/v1/external/upload",1)])])]),createBaseVNode("div",_hoisted_5$c,[createBaseVNode("button",{class:normalizeClass(["cyber-send-btn",{loading:__props.isLoading}]),onClick:_cache[0]||(_cache[0]=$event=>_ctx.$emit("send-request")),disabled:!__props.canSend||__props.isLoading},[createBaseVNode("div",_hoisted_7$c,[createBaseVNode("i",{class:normalizeClass(__props.isLoading?"fas fa-spinner fa-spin":"fas fa-paper-plane")},null,2),createBaseVNode("span",null,toDisplayString(__props.isLoading?unref($t2)("docs.tester.executing"):unref($t2)("docs.tester.send")),1)]),_cache[2]||(_cache[2]=createBaseVNode("div",{class:"btn-glow"},null,-1))],10,_hoisted_6$c)])]))}});const TesterHeader=_export_sfc(_sfc_main$i,[["__scopeId","data-v-4c2cf849"]]);function useFileUpload(selectedFiles,showValidation,isDragOver){const fileInput=ref(null);return{fileInput,formatFileSize:bytes=>{if(bytes===0)return"0 B";const k=1024,sizes=["B","KB","MB","GB"],i=Math.floor(Math.log(bytes)/Math.log(k));return parseFloat((bytes/Math.pow(k,i)).toFixed(1))+" "+sizes[i]},triggerFileInput:()=>{fileInput.value?.click()},handleFileSelect:event=>{const target=event.target;if(target.files){const files=Array.from(target.files);selectedFiles.value=[...selectedFiles.value,...files],showValidation.value=!1}},handleDrop:event=>{if(event.preventDefault(),isDragOver.value=!1,event.dataTransfer?.files){const imageFiles=Array.from(event.dataTransfer.files).filter(file=>file.type.startsWith("image/"));selectedFiles.value=[...selectedFiles.value,...imageFiles],showValidation.value=!1}},handleDragOver:()=>{isDragOver.value=!0},handleDragLeave:()=>{isDragOver.value=!1},removeFile:index2=>{selectedFiles.value.splice(index2,1)},handleImageError:event=>{const target=event.target;target.style.display="none"}}}const _hoisted_1$h={class:"cyber-section"},_hoisted_2$g={class:"section-header"},_hoisted_3$f={key:0,class:"upload-placeholder"},_hoisted_4$d={class:"upload-text"},_hoisted_5$b={class:"primary-text"},_hoisted_6$b={class:"secondary-text"},_hoisted_7$b={key:1,class:"files-list"},_hoisted_8$8={class:"files-header"},_hoisted_9$7={class:"files-count"},_hoisted_10$7={class:"files-size"},_hoisted_11$7={class:"files-container"},_hoisted_12$6={class:"file-info"},_hoisted_13$6={class:"file-name"},_hoisted_14$6={class:"file-size"},_hoisted_15$5=["onClick"],_sfc_main$h=defineComponent({name:"FileUpload",__name:"FileUpload",props:{selectedFiles:{},totalFileSize:{},showValidation:{type:Boolean},isDragOver:{type:Boolean}},emits:["update:selectedFiles","update:showValidation","update:isDragOver"],setup(__props,{emit:__emit}){const{$t:$t2}=useTexts(),props=__props,emit=__emit,selectedFiles=computed({get:()=>props.selectedFiles,set:value=>emit("update:selectedFiles",value)}),showValidation=computed({get:()=>props.showValidation,set:value=>emit("update:showValidation",value)}),isDragOver=computed({get:()=>props.isDragOver,set:value=>emit("update:isDragOver",value)}),{fileInput,formatFileSize,triggerFileInput,handleFileSelect,handleDrop,handleDragOver,handleDragLeave,removeFile}=useFileUpload(selectedFiles,showValidation,isDragOver);return(_ctx,_cache)=>(openBlock(),createElementBlock("div",_hoisted_1$h,[createBaseVNode("div",_hoisted_2$g,[_cache[6]||(_cache[6]=createBaseVNode("i",{class:"fas fa-file-image"},null,-1)),createBaseVNode("h4",null,toDisplayString(unref($t2)("docs.fileUpload.title")),1),_cache[7]||(_cache[7]=createBaseVNode("div",{class:"section-line"},null,-1))]),createBaseVNode("div",{class:normalizeClass(["cyber-upload-zone",{"has-files":selectedFiles.value.length>0,"drag-over":isDragOver.value,error:!selectedFiles.value.length&&showValidation.value}]),onClick:_cache[2]||(_cache[2]=(...args)=>unref(triggerFileInput)&&unref(triggerFileInput)(...args)),onDrop:_cache[3]||(_cache[3]=(...args)=>unref(handleDrop)&&unref(handleDrop)(...args)),onDragover:_cache[4]||(_cache[4]=withModifiers((...args)=>unref(handleDragOver)&&unref(handleDragOver)(...args),["prevent"])),onDragleave:_cache[5]||(_cache[5]=(...args)=>unref(handleDragLeave)&&unref(handleDragLeave)(...args))},[createBaseVNode("input",{ref_key:"fileInput",ref:fileInput,type:"file",multiple:"",accept:"image/*",onChange:_cache[0]||(_cache[0]=(...args)=>unref(handleFileSelect)&&unref(handleFileSelect)(...args)),style:{display:"none"}},null,544),selectedFiles.value.length?(openBlock(),createElementBlock("div",_hoisted_7$b,[createBaseVNode("div",_hoisted_8$8,[createBaseVNode("span",_hoisted_9$7,toDisplayString(unref($t2)("docs.fileUpload.selectedCount",{count:selectedFiles.value.length})),1),createBaseVNode("span",_hoisted_10$7,toDisplayString(unref($t2)("docs.fileUpload.totalSize",{size:unref(formatFileSize)(__props.totalFileSize)})),1)]),createBaseVNode("div",_hoisted_11$7,[(openBlock(!0),createElementBlock(Fragment,null,renderList(selectedFiles.value,(file,index2)=>(openBlock(),createElementBlock("div",{key:index2,class:"file-cyber-item"},[_cache[10]||(_cache[10]=createBaseVNode("div",{class:"file-icon"},[createBaseVNode("i",{class:"fas fa-file-image"})],-1)),createBaseVNode("div",_hoisted_12$6,[createBaseVNode("div",_hoisted_13$6,toDisplayString(file.name),1),createBaseVNode("div",_hoisted_14$6,toDisplayString(unref(formatFileSize)(file.size)),1)]),createBaseVNode("button",{onClick:withModifiers($event=>unref(removeFile)(index2),["stop"]),class:"remove-file-btn"},[..._cache[9]||(_cache[9]=[createBaseVNode("i",{class:"fas fa-times"},null,-1)])],8,_hoisted_15$5)]))),128))]),createBaseVNode("button",{onClick:_cache[1]||(_cache[1]=withModifiers((...args)=>unref(triggerFileInput)&&unref(triggerFileInput)(...args),["stop"])),class:"add-more-files-btn"},[_cache[11]||(_cache[11]=createBaseVNode("i",{class:"fas fa-plus"},null,-1)),createBaseVNode("span",null,toDisplayString(unref($t2)("docs.fileUpload.addMore")),1)])])):(openBlock(),createElementBlock("div",_hoisted_3$f,[_cache[8]||(_cache[8]=createBaseVNode("div",{class:"upload-icon"},[createBaseVNode("i",{class:"fas fa-images"})],-1)),createBaseVNode("div",_hoisted_4$d,[createBaseVNode("p",_hoisted_5$b,toDisplayString(unref($t2)("docs.fileUpload.dragOrClick")),1),createBaseVNode("p",_hoisted_6$b,toDisplayString(unref($t2)("docs.fileUpload.supportedFormats")),1)])]))],34)]))}});const FileUpload=_export_sfc(_sfc_main$h,[["__scopeId","data-v-e5820ad7"]]),_hoisted_1$g={class:"params-cyber-panel"},_hoisted_2$f={class:"cyber-section"},_hoisted_3$e={class:"section-header"},_hoisted_4$c={class:"auth-container"},_hoisted_5$a={class:"input-group"},_hoisted_6$a={class:"cyber-label"},_hoisted_7$a={class:"cyber-section"},_hoisted_8$7={class:"section-header"},_hoisted_9$6={class:"params-grid"},_hoisted_10$6={class:"input-group"},_hoisted_11$6={class:"cyber-label"},_hoisted_12$5={class:"input-group"},_hoisted_13$5={class:"cyber-label"},_hoisted_14$5={class:"input-group"},_hoisted_15$4={class:"cyber-label"},_hoisted_16$4={class:"input-group full-width"},_sfc_main$g=defineComponent({name:"ParamsPanel",__name:"ParamsPanel",props:{apiKey:{},selectedFiles:{},params:{},totalFileSize:{},showValidation:{type:Boolean},isDragOver:{type:Boolean}},emits:["update:apiKey","update:selectedFiles","update:params","update:showValidation","update:isDragOver"],setup(__props,{emit:__emit}){const props=__props,emit=__emit,{$t:$t2}=useTexts(),localSelectedFiles=computed({get:()=>props.selectedFiles,set:value=>emit("update:selectedFiles",value)}),localShowValidation=computed({get:()=>props.showValidation,set:value=>emit("update:showValidation",value)}),localIsDragOver=computed({get:()=>props.isDragOver,set:value=>emit("update:isDragOver",value)}),updateParam=(key,value)=>{const newParams={...props.params,[key]:value};emit("update:params",newParams)};return(_ctx,_cache)=>{const _component_CyberInput=resolveComponent("CyberInput"),_component_CyberRadioGroup=resolveComponent("CyberRadioGroup"),_component_CyberCheckbox=resolveComponent("CyberCheckbox");return openBlock(),createElementBlock("div",_hoisted_1$g,[createBaseVNode("div",_hoisted_2$f,[createBaseVNode("div",_hoisted_3$e,[_cache[8]||(_cache[8]=createBaseVNode("i",{class:"fas fa-shield-alt"},null,-1)),createBaseVNode("h4",null,toDisplayString(unref($t2)("docs.params.authConfig")),1),_cache[9]||(_cache[9]=createBaseVNode("div",{class:"section-line"},null,-1))]),createBaseVNode("div",_hoisted_4$c,[createBaseVNode("div",_hoisted_5$a,[createBaseVNode("label",_hoisted_6$a,[createBaseVNode("span",null,toDisplayString(unref($t2)("docs.params.apiKey")),1),_cache[10]||(_cache[10]=createBaseVNode("span",{class:"required-star"},"*",-1))]),createVNode(_component_CyberInput,{"model-value":__props.apiKey,"onUpdate:modelValue":_cache[0]||(_cache[0]=$event=>_ctx.$emit("update:apiKey",$event)),type:"password",placeholder:unref($t2)("docs.tester.apiKeyPlaceholder"),error:!__props.apiKey.trim()&&__props.showValidation},null,8,["model-value","placeholder","error"])])])]),createVNode(FileUpload,{selectedFiles:localSelectedFiles.value,"onUpdate:selectedFiles":_cache[1]||(_cache[1]=$event=>localSelectedFiles.value=$event),showValidation:localShowValidation.value,"onUpdate:showValidation":_cache[2]||(_cache[2]=$event=>localShowValidation.value=$event),isDragOver:localIsDragOver.value,"onUpdate:isDragOver":_cache[3]||(_cache[3]=$event=>localIsDragOver.value=$event),totalFileSize:__props.totalFileSize},null,8,["selectedFiles","showValidation","isDragOver","totalFileSize"]),createBaseVNode("div",_hoisted_7$a,[createBaseVNode("div",_hoisted_8$7,[_cache[11]||(_cache[11]=createBaseVNode("i",{class:"fas fa-cogs"},null,-1)),createBaseVNode("h4",null,toDisplayString(unref($t2)("docs.params.paramsConfig")),1),_cache[12]||(_cache[12]=createBaseVNode("div",{class:"section-line"},null,-1))]),createBaseVNode("div",_hoisted_9$6,[createBaseVNode("div",_hoisted_10$6,[createBaseVNode("label",_hoisted_11$6,toDisplayString(unref($t2)("docs.params.folderId")),1),createVNode(_component_CyberInput,{"model-value":__props.params.folderId,"onUpdate:modelValue":_cache[4]||(_cache[4]=$event=>updateParam("folderId",$event)),placeholder:unref($t2)("docs.tester.folderIdPlaceholder")},null,8,["model-value","placeholder"])]),createBaseVNode("div",_hoisted_12$5,[createBaseVNode("label",_hoisted_13$5,toDisplayString(unref($t2)("docs.params.filePath")),1),createVNode(_component_CyberInput,{"model-value":__props.params.filePath,"onUpdate:modelValue":_cache[5]||(_cache[5]=$event=>updateParam("filePath",$event)),placeholder:unref($t2)("docs.params.filePathPlaceholder")},null,8,["model-value","placeholder"])]),createBaseVNode("div",_hoisted_14$5,[createBaseVNode("label",_hoisted_15$4,toDisplayString(unref($t2)("docs.params.accessLevel")),1),createVNode(_component_CyberRadioGroup,{"model-value":__props.params.accessLevel,"onUpdate:modelValue":_cache[6]||(_cache[6]=$event=>updateParam("accessLevel",$event)),options:[{label:unref($t2)("docs.params.public"),value:"public"},{label:unref($t2)("docs.params.private"),value:"private"},{label:unref($t2)("docs.params.protected"),value:"protected"}]},null,8,["model-value","options"])]),createBaseVNode("div",_hoisted_16$4,[createVNode(_component_CyberCheckbox,{"model-value":__props.params.optimize,"onUpdate:modelValue":_cache[7]||(_cache[7]=$event=>updateParam("optimize",$event))},{default:withCtx(()=>[createTextVNode(toDisplayString(unref($t2)("docs.params.enableOptimize")),1)]),_:1},8,["model-value"])])])])])}}});const ParamsPanel=_export_sfc(_sfc_main$g,[["__scopeId","data-v-d9877dba"]]),_hoisted_1$f={class:"response-cyber-panel"},_hoisted_2$e={class:"section-header"},_hoisted_3$d={class:"response-content"},_hoisted_4$b={key:0,class:"empty-state"},_hoisted_5$9={class:"empty-container"},_hoisted_6$9={class:"empty-text-content"},_hoisted_7$9={class:"empty-title"},_hoisted_8$6={class:"empty-subtitle"},_hoisted_9$5={class:"test-guide"},_hoisted_10$5={class:"guide-steps"},_hoisted_11$5={class:"guide-step"},_hoisted_12$4={class:"step-text"},_hoisted_13$4={class:"guide-step"},_hoisted_14$4={class:"step-text"},_hoisted_15$3={class:"guide-step"},_hoisted_16$3={class:"step-text"},_hoisted_17$3={class:"guide-note"},_hoisted_18$3={key:1,class:"loading-state"},_hoisted_19$3={class:"loading-text"},_hoisted_20$3={key:2,class:"response-data"},_hoisted_21$2={class:"request-summary"},_hoisted_22$2={class:"summary-item"},_hoisted_23$2={class:"label"},_hoisted_24$2={class:"value"},_hoisted_25$2={class:"summary-item"},_hoisted_26$2={class:"label"},_hoisted_27$2={class:"value"},_hoisted_28$2={class:"summary-item"},_hoisted_29$2={class:"label"},_hoisted_30$1={class:"value"},_hoisted_31$1={class:"json-response"},_hoisted_32$1={class:"json-header"},_hoisted_33$1={class:"json-container"},_hoisted_34$1={class:"json-code"},_hoisted_35$1={key:0,class:"upload-results"},_hoisted_36$1={class:"results-title"},_hoisted_37={class:"images-grid"},_hoisted_38={class:"image-preview"},_hoisted_39=["src","alt"],_hoisted_40={class:"image-overlay"},_hoisted_41=["onClick","title"],_hoisted_42=["onClick","title"],_hoisted_43={class:"image-details"},_hoisted_44={class:"image-header"},_hoisted_45={class:"image-name"},_hoisted_46={class:"image-id"},_hoisted_47={class:"image-meta-grid"},_hoisted_48={class:"meta-item"},_hoisted_49={class:"meta-label"},_hoisted_50={class:"meta-value"},_hoisted_51={class:"meta-item"},_hoisted_52={class:"meta-label"},_hoisted_53={class:"meta-value"},_hoisted_54={class:"meta-item"},_hoisted_55={class:"meta-label"},_hoisted_56={class:"meta-value format-badge"},_hoisted_57={class:"meta-item"},_hoisted_58={class:"meta-label"},_hoisted_59={key:1,class:"error-results"},_hoisted_60={class:"results-title error"},_hoisted_61={key:0,class:"error-section"},_hoisted_62={class:"error-list"},_hoisted_63={class:"size-limit"},_hoisted_64={key:1,class:"error-section"},_hoisted_65={class:"error-list"},_hoisted_66={key:2,class:"error-section"},_hoisted_67={class:"error-list"},_hoisted_68={key:3,class:"error-section"},_hoisted_69={class:"error-list"},_sfc_main$f=defineComponent({name:"ResponsePanel",__name:"ResponsePanel",props:{response:{},isLoading:{type:Boolean},responseCopied:{type:Boolean},selectedFiles:{},totalFileSize:{},formattedResponse:{},responseStatusClass:{},uploadedImages:{},hasErrors:{type:Boolean},oversizedFiles:{},unsupportedFiles:{},invalidFiles:{},uploadErrors:{},sizeLimit:{},formatFileSize:{type:Function}},emits:["copy-response","copy-image-url","open-image"],setup(__props){const props=__props,{$t:$t2}=useTexts(),jsonCodeElement=ref(null);watch(()=>props.formattedResponse,()=>{nextTick(()=>{jsonCodeElement.value&&props.formattedResponse&&(delete jsonCodeElement.value.dataset.highlighted,HighlightJS.highlightElement(jsonCodeElement.value))})});const handleImageError=event=>{const target=event.target;target.style.display="none"};return(_ctx,_cache)=>(openBlock(),createElementBlock("div",_hoisted_1$f,[createBaseVNode("div",_hoisted_2$e,[_cache[2]||(_cache[2]=createBaseVNode("i",{class:"fas fa-terminal"},null,-1)),createBaseVNode("h4",null,toDisplayString(unref($t2)("docs.tester.response")),1),__props.response?(openBlock(),createElementBlock("div",{key:0,class:normalizeClass(["status-indicator",__props.responseStatusClass])},[_cache[1]||(_cache[1]=createBaseVNode("div",{class:"status-dot"},null,-1)),createBaseVNode("span",null,toDisplayString(__props.response.status)+" "+toDisplayString(__props.response.statusText),1)],2)):createCommentVNode("",!0),_cache[3]||(_cache[3]=createBaseVNode("div",{class:"section-line"},null,-1))]),createBaseVNode("div",_hoisted_3$d,[!__props.response&&!__props.isLoading?(openBlock(),createElementBlock("div",_hoisted_4$b,[createBaseVNode("div",_hoisted_5$9,[_cache[4]||(_cache[4]=createBaseVNode("div",{class:"empty-icon-wrapper"},[createBaseVNode("i",{class:"fas fa-flask"})],-1)),createBaseVNode("div",_hoisted_6$9,[createBaseVNode("h3",_hoisted_7$9,toDisplayString(unref($t2)("docs.response.emptyTitle")),1),createBaseVNode("p",_hoisted_8$6,toDisplayString(unref($t2)("docs.response.emptySubtitle")),1)])]),createBaseVNode("div",_hoisted_9$5,[createBaseVNode("div",_hoisted_10$5,[createBaseVNode("div",_hoisted_11$5,[_cache[5]||(_cache[5]=createBaseVNode("span",{class:"step-num"},"1",-1)),createBaseVNode("span",_hoisted_12$4,toDisplayString(unref($t2)("docs.response.step1")),1)]),_cache[8]||(_cache[8]=createBaseVNode("div",{class:"step-arrow"},"→",-1)),createBaseVNode("div",_hoisted_13$4,[_cache[6]||(_cache[6]=createBaseVNode("span",{class:"step-num"},"2",-1)),createBaseVNode("span",_hoisted_14$4,toDisplayString(unref($t2)("docs.response.step2")),1)]),_cache[9]||(_cache[9]=createBaseVNode("div",{class:"step-arrow"},"→",-1)),createBaseVNode("div",_hoisted_15$3,[_cache[7]||(_cache[7]=createBaseVNode("span",{class:"step-num"},"3",-1)),createBaseVNode("span",_hoisted_16$3,toDisplayString(unref($t2)("docs.response.step3")),1)])]),createBaseVNode("p",_hoisted_17$3,toDisplayString(unref($t2)("docs.tester.guide")),1)])])):__props.isLoading?(openBlock(),createElementBlock("div",_hoisted_18$3,[_cache[10]||(_cache[10]=createStaticVNode('<div class="loading-animation" data-v-0576636f><div class="loading-rings" data-v-0576636f><div class="ring ring-1" data-v-0576636f></div><div class="ring ring-2" data-v-0576636f></div><div class="ring-3 ring" data-v-0576636f></div></div></div>',1)),createBaseVNode("p",_hoisted_19$3,toDisplayString(unref($t2)("docs.response.processing")),1)])):(openBlock(),createElementBlock("div",_hoisted_20$3,[createBaseVNode("div",_hoisted_21$2,[createBaseVNode("div",_hoisted_22$2,[createBaseVNode("span",_hoisted_23$2,toDisplayString(unref($t2)("docs.response.fileCount")),1),createBaseVNode("span",_hoisted_24$2,toDisplayString(__props.selectedFiles.length),1)]),createBaseVNode("div",_hoisted_25$2,[createBaseVNode("span",_hoisted_26$2,toDisplayString(unref($t2)("docs.response.totalSize")),1),createBaseVNode("span",_hoisted_27$2,toDisplayString(__props.formatFileSize(__props.totalFileSize)),1)]),createBaseVNode("div",_hoisted_28$2,[createBaseVNode("span",_hoisted_29$2,toDisplayString(unref($t2)("docs.response.responseTime")),1),createBaseVNode("span",_hoisted_30$1,toDisplayString(__props.response.duration)+"ms",1)])]),createBaseVNode("div",_hoisted_31$1,[createBaseVNode("div",_hoisted_32$1,[createBaseVNode("span",null,toDisplayString(unref($t2)("docs.response.responseData")),1),createBaseVNode("button",{class:normalizeClass(["copy-json-btn",{copied:__props.responseCopied}]),onClick:_cache[0]||(_cache[0]=$event=>_ctx.$emit("copy-response"))},[createBaseVNode("i",{class:normalizeClass(__props.responseCopied?"fas fa-check":"fas fa-copy")},null,2),createBaseVNode("span",null,toDisplayString(__props.responseCopied?unref($t2)("docs.tester.copied"):unref($t2)("docs.tester.copy")),1)],2)]),createBaseVNode("div",_hoisted_33$1,[createBaseVNode("pre",_hoisted_34$1,[createBaseVNode("code",{ref_key:"jsonCodeElement",ref:jsonCodeElement,class:"language-json"},toDisplayString(__props.formattedResponse),513)])])]),__props.uploadedImages.length>0?(openBlock(),createElementBlock("div",_hoisted_35$1,[createBaseVNode("h5",_hoisted_36$1,[_cache[11]||(_cache[11]=createBaseVNode("i",{class:"fas fa-check-circle"},null,-1)),createTextVNode(" "+toDisplayString(unref($t2)("docs.response.uploadSuccess",{count:__props.uploadedImages.length})),1)]),createBaseVNode("div",_hoisted_37,[(openBlock(!0),createElementBlock(Fragment,null,renderList(__props.uploadedImages,image=>(openBlock(),createElementBlock("div",{key:image.id,class:"image-result-card"},[createBaseVNode("div",_hoisted_38,[createBaseVNode("img",{src:image.thumb_url,alt:image.original_name,onError:handleImageError},null,40,_hoisted_39),createBaseVNode("div",_hoisted_40,[createBaseVNode("button",{onClick:$event=>_ctx.$emit("copy-image-url",image.url),class:"copy-url-btn",title:unref($t2)("docs.tester.copyImageUrl")},[..._cache[12]||(_cache[12]=[createBaseVNode("i",{class:"fas fa-link"},null,-1)])],8,_hoisted_41),createBaseVNode("button",{onClick:$event=>_ctx.$emit("open-image",image.url),class:"view-original-btn",title:unref($t2)("docs.tester.viewOriginal")},[..._cache[13]||(_cache[13]=[createBaseVNode("i",{class:"fas fa-external-link-alt"},null,-1)])],8,_hoisted_42)])]),createBaseVNode("div",_hoisted_43,[createBaseVNode("div",_hoisted_44,[createBaseVNode("p",_hoisted_45,toDisplayString(image.display_name||image.original_name),1),createBaseVNode("span",_hoisted_46,toDisplayString(image.id),1)]),createBaseVNode("div",_hoisted_47,[createBaseVNode("div",_hoisted_48,[createBaseVNode("span",_hoisted_49,toDisplayString(unref($t2)("docs.response.dimensions")),1),createBaseVNode("span",_hoisted_50,toDisplayString(image.width)+"×"+toDisplayString(image.height),1)]),createBaseVNode("div",_hoisted_51,[createBaseVNode("span",_hoisted_52,toDisplayString(unref($t2)("docs.response.size")),1),createBaseVNode("span",_hoisted_53,toDisplayString(__props.formatFileSize(image.size)),1)]),createBaseVNode("div",_hoisted_54,[createBaseVNode("span",_hoisted_55,toDisplayString(unref($t2)("docs.response.format")),1),createBaseVNode("span",_hoisted_56,toDisplayString(image.format.toUpperCase()),1)]),createBaseVNode("div",_hoisted_57,[createBaseVNode("span",_hoisted_58,toDisplayString(unref($t2)("docs.response.accessLevel")),1),createBaseVNode("span",{class:normalizeClass(["meta-value access-level-badge",image.access_level])},toDisplayString(image.access_level==="public"?unref($t2)("docs.params.public"):image.access_level==="private"?unref($t2)("docs.params.private"):image.access_level==="protected"?unref($t2)("docs.params.protected"):image.access_level),3)])])])]))),128))])])):createCommentVNode("",!0),__props.hasErrors?(openBlock(),createElementBlock("div",_hoisted_59,[createBaseVNode("h5",_hoisted_60,[_cache[14]||(_cache[14]=createBaseVNode("i",{class:"fas fa-exclamation-triangle"},null,-1)),createTextVNode(" "+toDisplayString(unref($t2)("docs.response.processingError")),1)]),__props.oversizedFiles.length>0?(openBlock(),createElementBlock("div",_hoisted_61,[createBaseVNode("h6",null,toDisplayString(unref($t2)("docs.response.oversizedFiles")),1),createBaseVNode("ul",_hoisted_62,[(openBlock(!0),createElementBlock(Fragment,null,renderList(__props.oversizedFiles,file=>(openBlock(),createElementBlock("li",{key:file},toDisplayString(file),1))),128))]),createBaseVNode("p",_hoisted_63,toDisplayString(unref($t2)("docs.response.sizeLimit",{limit:__props.sizeLimit})),1)])):createCommentVNode("",!0),__props.unsupportedFiles.length>0?(openBlock(),createElementBlock("div",_hoisted_64,[createBaseVNode("h6",null,toDisplayString(unref($t2)("docs.response.unsupportedFiles")),1),createBaseVNode("ul",_hoisted_65,[(openBlock(!0),createElementBlock(Fragment,null,renderList(__props.unsupportedFiles,file=>(openBlock(),createElementBlock("li",{key:file},toDisplayString(file),1))),128))])])):createCommentVNode("",!0),__props.invalidFiles.length>0?(openBlock(),createElementBlock("div",_hoisted_66,[createBaseVNode("h6",null,toDisplayString(unref($t2)("docs.response.invalidFiles")),1),createBaseVNode("ul",_hoisted_67,[(openBlock(!0),createElementBlock(Fragment,null,renderList(__props.invalidFiles,file=>(openBlock(),createElementBlock("li",{key:file},toDisplayString(file),1))),128))])])):createCommentVNode("",!0),__props.uploadErrors.length>0?(openBlock(),createElementBlock("div",_hoisted_68,[createBaseVNode("h6",null,toDisplayString(unref($t2)("docs.response.uploadFailed")),1),createBaseVNode("ul",_hoisted_69,[(openBlock(!0),createElementBlock(Fragment,null,renderList(__props.uploadErrors,error=>(openBlock(),createElementBlock("li",{key:error},toDisplayString(error),1))),128))])])):createCommentVNode("",!0)])):createCommentVNode("",!0)]))])]))}});const ResponsePanel=_export_sfc(_sfc_main$f,[["__scopeId","data-v-0576636f"]]);function useApiTest(){const toast=useToast(),{$t:$t2}=useTexts(),apiKey=ref(""),selectedFiles=ref([]),params=ref({folderId:"",filePath:"",accessLevel:"private",optimize:!1}),response=ref(null),isLoading=ref(!1),responseCopied=ref(!1),showValidation=ref(!1),isDragOver=ref(!1),currentDomain=computed(()=>{const apiBaseUrl=API_BASE_URL;return apiBaseUrl.includes("/api")?apiBaseUrl.replace(/\/api.*$/,""):apiBaseUrl.startsWith("/")?typeof window<"u"?window.location.origin:"http://localhost:9520":apiBaseUrl}),totalFileSize=computed(()=>selectedFiles.value.reduce((total,file)=>total+file.size,0)),responseStatusClass=computed(()=>{if(!response.value)return"";const status=response.value.status;return status>=200&&status<300?"success":status>=400&&status<500?"error":status>=500?"server-error":""}),uploadedImages=computed(()=>response.value?.data?.data?.uploaded||[]),hasErrors=computed(()=>{const data=response.value?.data?.data;return!!(data?.oversized_files&&data.oversized_files.length>0||data?.unsupported_files&&data.unsupported_files.length>0||data?.invalid_files&&data.invalid_files.length>0||data?.upload_errors&&data.upload_errors.length>0)}),oversizedFiles=computed(()=>response.value?.data?.data?.oversized_files||[]),unsupportedFiles=computed(()=>response.value?.data?.data?.unsupported_files||[]),invalidFiles=computed(()=>response.value?.data?.data?.invalid_files||[]),uploadErrors=computed(()=>response.value?.data?.data?.upload_errors||[]),sizeLimit=computed(()=>response.value?.data?.data?.size_limit||"");return{apiKey,selectedFiles,params,response,isLoading,responseCopied,showValidation,isDragOver,currentDomain,totalFileSize,responseStatusClass,uploadedImages,hasErrors,oversizedFiles,unsupportedFiles,invalidFiles,uploadErrors,sizeLimit,sendRequest:async()=>{if(!selectedFiles.value.length||!apiKey.value.trim()){showValidation.value=!0;return}isLoading.value=!0,response.value=null,responseCopied.value=!1;try{const formData=new FormData;selectedFiles.value.forEach(file=>{formData.append("files[]",file)}),params.value.folderId&&formData.append("folderId",params.value.folderId),params.value.filePath&&formData.append("filePath",params.value.filePath),formData.append("access_level",params.value.accessLevel),formData.append("optimize",params.value.optimize?"true":"false");const startTime=Date.now(),res=await fetch(`${currentDomain.value}/api/v1/external/upload`,{method:"POST",headers:{"x-pixelpunk-key":apiKey.value},body:formData}),endTime=Date.now(),data=await res.json();response.value={status:res.status,statusText:res.statusText,data,duration:endTime-startTime}}catch(error){response.value={status:0,statusText:"Network Error",data:{error:error.message},duration:0}}finally{isLoading.value=!1}},copyResponse:async()=>{try{const jsonString=JSON.stringify(response.value.data,null,2);await navigator.clipboard.writeText(jsonString),responseCopied.value=!0,toast.success($t2("docs.apiTester.toast.responseCopied")),setTimeout(()=>{responseCopied.value=!1},2e3)}catch(error){console.error("Copy failed:",error);try{const jsonString=JSON.stringify(response.value.data,null,2),textArea=document.createElement("textarea");textArea.value=jsonString,textArea.style.position="fixed",textArea.style.left="-999999px",textArea.style.top="-999999px",document.body.appendChild(textArea),textArea.focus(),textArea.select();const success=document.execCommand("copy");if(document.body.removeChild(textArea),success)responseCopied.value=!0,toast.success($t2("docs.apiTester.toast.responseCopied")),setTimeout(()=>{responseCopied.value=!1},2e3);else throw new Error("execCommand failed")}catch(fallbackError){console.error("fallback复制也失败:",fallbackError),toast.error($t2("docs.apiTester.toast.copyFailed"))}}},copyImageUrl:async url=>{try{await navigator.clipboard.writeText(url),toast.success($t2("docs.apiTester.toast.imageLinkCopied"))}catch(error){console.error("Copy failed:",error);try{const textArea=document.createElement("textarea");textArea.value=url,textArea.style.position="fixed",textArea.style.left="-999999px",textArea.style.top="-999999px",document.body.appendChild(textArea),textArea.focus(),textArea.select();const success=document.execCommand("copy");if(document.body.removeChild(textArea),success)toast.success($t2("docs.apiTester.toast.imageLinkCopied"));else throw new Error("execCommand failed")}catch(fallbackError){console.error("fallback复制也失败:",fallbackError),toast.error($t2("docs.apiTester.toast.copyLinkFailed"))}}},openImageInNewTab:url=>{try{window.open(url,"_blank")?toast.success($t2("docs.apiTester.toast.openingNewWindow")):toast.error($t2("docs.apiTester.toast.unableToOpenWindow"))}catch(error){console.error("Failed to open new window:",error),toast.error($t2("docs.apiTester.toast.failedToOpenWindow"))}}}}function useJsonFormatter(response){return{formattedResponse:computed(()=>response.value?.data?JSON.stringify(response.value.data,null,2):""),formatDate:dateString=>dateString?new Date(dateString).toLocaleDateString(getCurrentLocale(),{year:"numeric",month:"2-digit",day:"2-digit",hour:"2-digit",minute:"2-digit"}):""}}const _hoisted_1$e={class:"api-tester-cyber"},_hoisted_2$d={class:"cyber-tester-body"},_sfc_main$e=defineComponent({name:"ApiTester",__name:"ApiTester",setup(__props){const apiTestComposable=useApiTest(),{apiKey,selectedFiles,params,response,isLoading,responseCopied,showValidation,isDragOver,currentDomain,totalFileSize,responseStatusClass,uploadedImages,hasErrors,oversizedFiles,unsupportedFiles,invalidFiles,uploadErrors,sizeLimit,sendRequest,copyResponse,copyImageUrl,openImageInNewTab}=apiTestComposable,{formatFileSize}=useFileUpload(selectedFiles,showValidation,isDragOver),{formattedResponse}=useJsonFormatter(response);return(_ctx,_cache)=>(openBlock(),createElementBlock("div",_hoisted_1$e,[createVNode(TesterHeader,{currentDomain:unref(currentDomain),canSend:unref(selectedFiles).length>0&&unref(apiKey).trim().length>0,isLoading:unref(isLoading),onSendRequest:unref(sendRequest)},null,8,["currentDomain","canSend","isLoading","onSendRequest"]),createBaseVNode("div",_hoisted_2$d,[createVNode(ParamsPanel,{apiKey:unref(apiKey),"onUpdate:apiKey":_cache[0]||(_cache[0]=$event=>isRef(apiKey)?apiKey.value=$event:null),selectedFiles:unref(selectedFiles),"onUpdate:selectedFiles":_cache[1]||(_cache[1]=$event=>isRef(selectedFiles)?selectedFiles.value=$event:null),params:unref(params),"onUpdate:params":_cache[2]||(_cache[2]=$event=>isRef(params)?params.value=$event:null),showValidation:unref(showValidation),"onUpdate:showValidation":_cache[3]||(_cache[3]=$event=>isRef(showValidation)?showValidation.value=$event:null),isDragOver:unref(isDragOver),"onUpdate:isDragOver":_cache[4]||(_cache[4]=$event=>isRef(isDragOver)?isDragOver.value=$event:null),totalFileSize:unref(totalFileSize)},null,8,["apiKey","selectedFiles","params","showValidation","isDragOver","totalFileSize"]),createVNode(ResponsePanel,{response:unref(response),isLoading:unref(isLoading),responseCopied:unref(responseCopied),selectedFiles:unref(selectedFiles),totalFileSize:unref(totalFileSize),formattedResponse:unref(formattedResponse),responseStatusClass:unref(responseStatusClass),uploadedImages:unref(uploadedImages),hasErrors:unref(hasErrors),oversizedFiles:unref(oversizedFiles),unsupportedFiles:unref(unsupportedFiles),invalidFiles:unref(invalidFiles),uploadErrors:unref(uploadErrors),sizeLimit:unref(sizeLimit),formatFileSize:unref(formatFileSize),onCopyResponse:unref(copyResponse),onCopyImageUrl:unref(copyImageUrl),onOpenImage:unref(openImageInNewTab)},null,8,["response","isLoading","responseCopied","selectedFiles","totalFileSize","formattedResponse","responseStatusClass","uploadedImages","hasErrors","oversizedFiles","unsupportedFiles","invalidFiles","uploadErrors","sizeLimit","formatFileSize","onCopyResponse","onCopyImageUrl","onOpenImage"])])]))}});const ApiTester=_export_sfc(_sfc_main$e,[["__scopeId","data-v-33d1fac3"]]),_hoisted_1$d={class:"docs-intro"},_hoisted_2$c={class:"title-section"},_hoisted_3$c={class:"title-wrapper"},_hoisted_4$a={class:"page-title"},_hoisted_5$8={class:"title-text"},_hoisted_6$8={class:"page-subtitle"},_hoisted_7$8={class:"status-indicators"},_hoisted_8$5={class:"status-item"},_hoisted_9$4={class:"status-item"},_hoisted_10$4={class:"status-item"},_hoisted_11$4={class:"quick-actions"},_sfc_main$d=defineComponent({name:"DocsIntro",__name:"index",props:{scrollToSection:{type:Function}},setup(__props){const{$t:$t2}=useTexts(),props=__props,onScrollToSection=target=>{props.scrollToSection(target)};return(_ctx,_cache)=>{const _component_CyberButton=resolveComponent("CyberButton");return openBlock(),createElementBlock("div",_hoisted_1$d,[_cache[11]||(_cache[11]=createBaseVNode("div",{class:"intro-background"},[createBaseVNode("div",{class:"intro-grid"}),createBaseVNode("div",{class:"intro-particles"})],-1)),createBaseVNode("div",_hoisted_2$c,[createBaseVNode("div",_hoisted_3$c,[createBaseVNode("h1",_hoisted_4$a,[_cache[3]||(_cache[3]=createBaseVNode("i",{class:"fas fa-code title-icon"},null,-1)),createBaseVNode("span",_hoisted_5$8,toDisplayString(unref($t2)("docs.intro.pageTitle")),1),_cache[4]||(_cache[4]=createBaseVNode("div",{class:"title-glow"},null,-1))]),createBaseVNode("p",_hoisted_6$8,toDisplayString(unref($t2)("docs.intro.subtitle")),1)]),createBaseVNode("div",_hoisted_7$8,[createBaseVNode("div",_hoisted_8$5,[_cache[5]||(_cache[5]=createBaseVNode("div",{class:"status-dot online"},null,-1)),createBaseVNode("span",null,toDisplayString(unref($t2)("docs.intro.status.online")),1)]),createBaseVNode("div",_hoisted_9$4,[_cache[6]||(_cache[6]=createBaseVNode("div",{class:"status-dot"},null,-1)),createBaseVNode("span",null,toDisplayString(unref($t2)("docs.intro.status.version")),1)]),createBaseVNode("div",_hoisted_10$4,[_cache[7]||(_cache[7]=createBaseVNode("div",{class:"status-dot perf"},null,-1)),createBaseVNode("span",null,toDisplayString(unref($t2)("docs.intro.status.detailed")),1)])])]),createBaseVNode("div",_hoisted_11$4,[createVNode(_component_CyberButton,{type:"primary",onClick:_cache[0]||(_cache[0]=$event=>onScrollToSection("#api-examples"))},{default:withCtx(()=>[_cache[8]||(_cache[8]=createBaseVNode("i",{class:"fas fa-code"},null,-1)),createBaseVNode("span",null,toDisplayString(unref($t2)("docs.intro.actions.examples")),1)]),_:1}),createVNode(_component_CyberButton,{type:"success",onClick:_cache[1]||(_cache[1]=$event=>onScrollToSection("#api-upload"))},{default:withCtx(()=>[_cache[9]||(_cache[9]=createBaseVNode("i",{class:"fas fa-upload"},null,-1)),createBaseVNode("span",null,toDisplayString(unref($t2)("docs.intro.actions.start")),1)]),_:1}),createVNode(_component_CyberButton,{type:"warning",onClick:_cache[2]||(_cache[2]=$event=>onScrollToSection("#api-tester"))},{default:withCtx(()=>[_cache[10]||(_cache[10]=createBaseVNode("i",{class:"fas fa-flask"},null,-1)),createBaseVNode("span",null,toDisplayString(unref($t2)("docs.intro.actions.test")),1)]),_:1})])])}}});const DocsIntro=_export_sfc(_sfc_main$d,[["__scopeId","data-v-ddd40eb8"]]),_hoisted_1$c={class:"docs-sidebar"},_hoisted_2$b={class:"sidebar-content"},_hoisted_3$b={class:"sidebar-header"},_hoisted_4$9={class:"brand-section"},_hoisted_5$7={class:"brand-text"},_hoisted_6$7={class:"brand-title"},_hoisted_7$7={class:"brand-subtitle"},_hoisted_8$4={class:"sidebar-nav-section"},_hoisted_9$3={class:"nav-header"},_hoisted_10$3={class:"nav-title-row"},_hoisted_11$3={class:"progress-percentage"},_hoisted_12$3={class:"progress-container"},_hoisted_13$3={class:"progress-bar"},_hoisted_14$3=["data-progress"],_hoisted_15$2={class:"sidebar-nav"},_hoisted_16$2=["href"],_hoisted_17$2={class:"link-content"},_hoisted_18$2={class:"link-icon-wrapper"},_hoisted_19$2={class:"link-text"},_hoisted_20$2={key:0,class:"active-indicator"},_sfc_main$c=defineComponent({name:"DocsSidebar",__name:"index",props:{navItems:{},activeSection:{},readingProgress:{},scrollToSection:{type:Function}},setup(__props){const{$t:$t2}=useTexts(),props=__props,progressStyle=computed(()=>({width:`${Math.max(0,Math.min(100,props.readingProgress||0))}%`})),onScrollToSection=event=>{props.scrollToSection(event)};return(_ctx,_cache)=>{const _component_router_link=resolveComponent("router-link");return openBlock(),createElementBlock("aside",_hoisted_1$c,[createBaseVNode("div",_hoisted_2$b,[createBaseVNode("div",_hoisted_3$b,[createVNode(_component_router_link,{to:"/",class:"header-brand-link"},{default:withCtx(()=>[createBaseVNode("div",_hoisted_4$9,[_cache[0]||(_cache[0]=createBaseVNode("div",{class:"brand-icon-wrapper"},[createBaseVNode("i",{class:"fas fa-arrow-left back-icon"}),createBaseVNode("i",{class:"fas fa-code main-icon"})],-1)),createBaseVNode("div",_hoisted_5$7,[createBaseVNode("span",_hoisted_6$7,toDisplayString(unref($t2)("docs.sidebar.title")),1),createBaseVNode("span",_hoisted_7$7,toDisplayString(unref($t2)("docs.sidebar.backHome")),1)])])]),_:1})]),createBaseVNode("div",_hoisted_8$4,[createBaseVNode("div",_hoisted_9$3,[createBaseVNode("div",_hoisted_10$3,[createBaseVNode("h3",null,[_cache[1]||(_cache[1]=createBaseVNode("i",{class:"fas fa-list-ul"},null,-1)),createTextVNode(toDisplayString(unref($t2)("docs.sidebar.navigation")),1)]),createBaseVNode("div",_hoisted_11$3,toDisplayString(Math.round(__props.readingProgress))+"%",1)]),createBaseVNode("div",_hoisted_12$3,[createBaseVNode("div",_hoisted_13$3,[createBaseVNode("div",{class:"progress-fill",style:normalizeStyle(progressStyle.value),"data-progress":__props.readingProgress},null,12,_hoisted_14$3)])])]),createBaseVNode("nav",_hoisted_15$2,[(openBlock(!0),createElementBlock(Fragment,null,renderList(__props.navItems,item=>(openBlock(),createElementBlock("a",{key:item.id,href:`#${item.id}`,class:normalizeClass(["sidebar-link",[{active:__props.activeSection===item.id}]]),onClick:onScrollToSection},[createBaseVNode("div",_hoisted_17$2,[createBaseVNode("div",_hoisted_18$2,[createBaseVNode("i",{class:normalizeClass(item.iconClass)},null,2)]),createBaseVNode("span",_hoisted_19$2,toDisplayString(item.label),1)]),__props.activeSection===item.id?(openBlock(),createElementBlock("div",_hoisted_20$2)):createCommentVNode("",!0)],10,_hoisted_16$2))),128))]),_cache[2]||(_cache[2]=createBaseVNode("div",{class:"sidebar-footer"},[createBaseVNode("div",{class:"api-info"},[createBaseVNode("div",{class:"info-item"},[createBaseVNode("i",{class:"fas fa-shield-alt"}),createBaseVNode("span",null,"API v1.0")])])],-1))])])])}}});const DocsSidebar=_export_sfc(_sfc_main$c,[["__scopeId","data-v-67eb37c6"]]),_hoisted_1$b={class:"mobile-header mobile-only"},_hoisted_2$a={class:"mobile-header-content"},_hoisted_3$a={class:"header-brand"},_hoisted_4$8=["href"],_sfc_main$b=defineComponent({name:"MobileHeader",__name:"index",props:{navItems:{},activeSection:{},mobileMenuOpen:{type:Boolean},toggleMobileMenu:{type:Function},handleMobileNavClick:{type:Function}},setup(__props){const{$t:$t2}=useTexts(),props=__props,onToggleMobileMenu=()=>{props.toggleMobileMenu()},onMobileNavClick=event=>{props.handleMobileNavClick(event)};return(_ctx,_cache)=>{const _component_router_link=resolveComponent("router-link");return openBlock(),createElementBlock("div",_hoisted_1$b,[createBaseVNode("div",_hoisted_2$a,[createVNode(_component_router_link,{to:"/",class:"back-home-btn"},{default:withCtx(()=>[_cache[0]||(_cache[0]=createBaseVNode("i",{class:"fas fa-arrow-left"},null,-1)),createBaseVNode("span",null,toDisplayString(unref($t2)("docs.mobile.backHome")),1)]),_:1}),createBaseVNode("div",_hoisted_3$a,[_cache[1]||(_cache[1]=createBaseVNode("i",{class:"fas fa-code"},null,-1)),createBaseVNode("span",null,toDisplayString(unref($t2)("docs.mobile.title")),1)]),createBaseVNode("button",{class:"mobile-menu-toggle",onClick:onToggleMobileMenu},[createBaseVNode("i",{class:normalizeClass(__props.mobileMenuOpen?"fas fa-times":"fas fa-bars")},null,2)])]),createBaseVNode("div",{class:normalizeClass(["mobile-menu",[{open:__props.mobileMenuOpen}]])},[(openBlock(!0),createElementBlock(Fragment,null,renderList(__props.navItems,item=>(openBlock(),createElementBlock("a",{key:item.id,href:`#${item.id}`,class:normalizeClass(["mobile-nav-link",[{active:__props.activeSection===item.id}]]),onClick:onMobileNavClick},[createBaseVNode("i",{class:normalizeClass(item.iconClass)},null,2),createBaseVNode("span",null,toDisplayString(item.label),1)],10,_hoisted_4$8))),128))],2)])}}});const MobileHeader=_export_sfc(_sfc_main$b,[["__scopeId","data-v-df4fe0c5"]]),_hoisted_1$a=["title"],_sfc_main$a=defineComponent({name:"BackToTop",__name:"index",props:{showBackToTop:{type:Boolean},backToTop:{type:Function}},setup(__props){const{$t:$t2}=useTexts(),props=__props,onBackToTop=()=>{props.backToTop()};return(_ctx,_cache)=>withDirectives((openBlock(),createElementBlock("button",{class:"back-to-top",title:unref($t2)("docs.backToTop"),onClick:onBackToTop},[..._cache[0]||(_cache[0]=[createBaseVNode("i",{class:"fas fa-arrow-up"},null,-1)])],8,_hoisted_1$a)),[[vShow,__props.showBackToTop]])}});const BackToTop=_export_sfc(_sfc_main$a,[["__scopeId","data-v-29d355f9"]]),_sfc_main$9=defineComponent({name:"ApiOverview",components:{ApiSection},props:{currentDomain:{type:String,required:!0}},setup(){const{$t:$t2}=useTexts();return{$t:$t2}}});const _hoisted_1$9={class:"section-intro"},_hoisted_2$9={class:"info-cards"},_hoisted_3$9={class:"info-card"},_hoisted_4$7={class:"card-content"},_hoisted_5$6={class:"info-card"},_hoisted_6$6={class:"card-content"},_hoisted_7$6={class:"info-card"},_hoisted_8$3={class:"card-content"};function _sfc_render$4(_ctx,_cache,$props,$setup,$data,$options){const _component_ApiSection=resolveComponent("ApiSection");return openBlock(),createBlock(_component_ApiSection,{id:"api-overview",title:_ctx.$t("docs.overview.title"),"icon-class":"fas fa-book"},{default:withCtx(()=>[createBaseVNode("div",_hoisted_1$9,[createBaseVNode("p",null,toDisplayString(_ctx.$t("docs.overview.description")),1)]),createBaseVNode("div",_hoisted_2$9,[createBaseVNode("div",_hoisted_3$9,[_cache[0]||(_cache[0]=createBaseVNode("div",{class:"card-icon"},[createBaseVNode("i",{class:"fas fa-server"})],-1)),createBaseVNode("div",_hoisted_4$7,[createBaseVNode("h4",null,toDisplayString(_ctx.$t("docs.overview.apiAddress")),1),createBaseVNode("code",null,toDisplayString(_ctx.currentDomain)+"/api/v1",1)])]),createBaseVNode("div",_hoisted_5$6,[_cache[1]||(_cache[1]=createBaseVNode("div",{class:"card-icon"},[createBaseVNode("i",{class:"fas fa-shield-alt"})],-1)),createBaseVNode("div",_hoisted_6$6,[createBaseVNode("h4",null,toDisplayString(_ctx.$t("docs.overview.authMethod")),1),createBaseVNode("p",null,toDisplayString(_ctx.$t("docs.overview.apiKeyAuth")),1)])]),createBaseVNode("div",_hoisted_7$6,[_cache[2]||(_cache[2]=createBaseVNode("div",{class:"card-icon"},[createBaseVNode("i",{class:"fas fa-zap"})],-1)),createBaseVNode("div",_hoisted_8$3,[createBaseVNode("h4",null,toDisplayString(_ctx.$t("docs.overview.performance")),1),createBaseVNode("p",null,toDisplayString(_ctx.$t("docs.overview.avgResponse")),1)])])])]),_:1},8,["title"])}const ApiOverview=_export_sfc(_sfc_main$9,[["render",_sfc_render$4],["__scopeId","data-v-b534d038"]]),getApiTesterFeatures=()=>{const{$t:$t2}=useTexts();return[{id:"file-upload",icon:"fas fa-upload",title:$t2("docs.apiTester.features.fileUpload.title"),description:$t2("docs.apiTester.features.fileUpload.description")},{id:"params-config",icon:"fas fa-cogs",title:$t2("docs.apiTester.features.paramsConfig.title"),description:$t2("docs.apiTester.features.paramsConfig.description")},{id:"real-time-response",icon:"fas fa-eye",title:$t2("docs.apiTester.features.realTimeResponse.title"),description:$t2("docs.apiTester.features.realTimeResponse.description")},{id:"result-preview",icon:"fas fa-image",title:$t2("docs.apiTester.features.resultPreview.title"),description:$t2("docs.apiTester.features.resultPreview.description")}]};getApiTesterFeatures();const BACK_TO_TOP_THRESHOLD=300,SCROLL_OFFSET=100,SECTION_ACTIVE_OFFSET=150,COPY_SUCCESS_DURATION=2e3,DOM_INIT_DELAY=100,SCROLL_REINIT_DELAY=200,{$t}=useTexts(),navItems=computed(()=>getDocsNavItems($t)),codeTabs=CODE_TAB_OPTIONS,supportedFormats=getImageFormatOptions($t).map(option=>option.value.toUpperCase());function useDocsPage(){const activeSection=ref("api-overview"),activeCodeTab=ref("curl"),copyStatus=ref({}),mobileMenuOpen=ref(!1),showBackToTop=ref(!1),readingProgress=ref(0),currentDomain=computed(()=>typeof window<"u"?window.location.origin:"https://your-domain.com"),toggleMobileMenu=()=>{mobileMenuOpen.value=!mobileMenuOpen.value},handleMobileNavClick=event=>{scrollToSection(event),mobileMenuOpen.value=!1},scrollToSection=eventOrTarget=>{eventOrTarget&&typeof eventOrTarget=="object"&&eventOrTarget.preventDefault&&eventOrTarget.preventDefault();let targetId="";if(typeof eventOrTarget=="string"?targetId=eventOrTarget.replace("#",""):eventOrTarget&&eventOrTarget.target&&(targetId=eventOrTarget.target.getAttribute("href")?.substring(1)||eventOrTarget.currentTarget.getAttribute("href")?.substring(1)),targetId){const element=document.getElementById(targetId);if(element){const docsContent=document.querySelector(".docs-content");if(docsContent){const elementTop=element.offsetTop;docsContent.scrollTo({top:elementTop-SCROLL_OFFSET,behavior:"smooth"})}activeSection.value=targetId}}},backToTop=()=>{const docsContent=document.querySelector(".docs-content");docsContent&&docsContent.scrollTo({top:0,behavior:"smooth"})},handleScroll=()=>{const docsContent=document.querySelector(".docs-content");if(!docsContent)return;const{scrollTop,scrollHeight,clientHeight}=docsContent,docHeight=scrollHeight-clientHeight;if(docHeight>0){const progress=Math.min(scrollTop/docHeight*100,100);readingProgress.value=progress}else readingProgress.value=0;showBackToTop.value=scrollTop>BACK_TO_TOP_THRESHOLD,updateActiveSection()},updateActiveSection=()=>{const docsContent=document.querySelector(".docs-content");if(!docsContent)return;const sections=document.querySelectorAll('[id^="api-"]'),{scrollTop}=docsContent;let current="";sections.forEach(section=>{const offsetTop=section.offsetTop-SECTION_ACTIVE_OFFSET;scrollTop>=offsetTop&&(current=section.id)}),current&&current!==activeSection.value&&(activeSection.value=current)},copyCode=(codeId,code)=>{code&&navigator.clipboard.writeText(code).then(()=>{copyStatus.value[codeId]=!0,setTimeout(()=>{copyStatus.value[codeId]=!1},COPY_SUCCESS_DURATION)}).catch(err=>{logger.error("Copy failed:",err)})},switchCodeTab=tabId=>{activeCodeTab.value=tabId,nextTick(()=>{highlightCodeBlocks()})},highlightCodeBlocks=()=>{document.querySelectorAll("pre code:not([data-highlighted])").forEach(block=>{HighlightJS.highlightElement(block)})},handleClickOutside=event=>{const mobileMenu=document.querySelector(".mobile-menu"),menuBtn=document.querySelector(".mobile-menu-btn");mobileMenuOpen.value&&mobileMenu&&!mobileMenu.contains(event.target)&&!menuBtn.contains(event.target)&&(mobileMenuOpen.value=!1)};return onMounted(()=>{setTimeout(()=>{const docsContent=document.querySelector(".docs-content");docsContent&&(docsContent.addEventListener("scroll",handleScroll,{passive:!0}),handleScroll())},DOM_INIT_DELAY),document.addEventListener("click",handleClickOutside),nextTick(()=>{highlightCodeBlocks(),setTimeout(()=>{handleScroll()},SCROLL_REINIT_DELAY)})}),onUnmounted(()=>{const docsContent=document.querySelector(".docs-content");docsContent&&docsContent.removeEventListener("scroll",handleScroll),document.removeEventListener("click",handleClickOutside)}),{activeSection,activeCodeTab,copyStatus,mobileMenuOpen,showBackToTop,readingProgress,currentDomain,toggleMobileMenu,handleMobileNavClick,scrollToSection,backToTop,copyCode,switchCodeTab,highlightCodeBlocks}}const _hoisted_1$8={class:"code-tabs"},_hoisted_2$8={class:"code-tabs-header"},_hoisted_3$8=["onClick"],_hoisted_4$6={class:"code-tabs-content"},_sfc_main$8=defineComponent({name:"CodeTabs",__name:"index",props:{tabs:{},default:{default:"curl"}},emits:["update:tab"],setup(__props,{emit:__emit}){const props=__props,emit=__emit,activeTab=ref(props.default);function changeTab(tab){activeTab.value=tab,emit("update:tab",tab)}return watch(()=>props.default,newValue=>{activeTab.value=newValue}),(_ctx,_cache)=>(openBlock(),createElementBlock("div",_hoisted_1$8,[createBaseVNode("div",_hoisted_2$8,[(openBlock(!0),createElementBlock(Fragment,null,renderList(__props.tabs,tab=>(openBlock(),createElementBlock("button",{key:tab.value,class:normalizeClass(["tab-button",[{active:activeTab.value===tab.value}]]),onClick:$event=>changeTab(tab.value)},toDisplayString(tab.label),11,_hoisted_3$8))),128))]),createBaseVNode("div",_hoisted_4$6,[activeTab.value==="curl"?renderSlot(_ctx.$slots,"curl",{key:0},void 0,!0):createCommentVNode("",!0),activeTab.value==="python"?renderSlot(_ctx.$slots,"python",{key:1},void 0,!0):createCommentVNode("",!0),activeTab.value==="nodejs"?renderSlot(_ctx.$slots,"nodejs",{key:2},void 0,!0):createCommentVNode("",!0),activeTab.value==="javascript"?renderSlot(_ctx.$slots,"javascript",{key:3},void 0,!0):createCommentVNode("",!0),activeTab.value==="java"?renderSlot(_ctx.$slots,"java",{key:4},void 0,!0):createCommentVNode("",!0),activeTab.value==="go"?renderSlot(_ctx.$slots,"go",{key:5},void 0,!0):createCommentVNode("",!0)])]))}});const CodeTabs=_export_sfc(_sfc_main$8,[["__scopeId","data-v-58dfcc75"]]),_hoisted_1$7={class:"code-example"},_hoisted_2$7={class:"code-header"},_hoisted_3$7={class:"code-title"},_hoisted_4$5={class:"code-container"},_hoisted_5$5={class:"code-block"},_hoisted_6$5={class:"language-bash"},_hoisted_7$5={class:"code-example"},_hoisted_8$2={class:"code-header"},_hoisted_9$2={class:"code-title"},_hoisted_10$2={class:"code-container"},_hoisted_11$2={class:"code-block"},_hoisted_12$2={class:"language-javascript"},_hoisted_13$2={class:"code-example"},_hoisted_14$2={class:"code-header"},_hoisted_15$1={class:"code-title"},_hoisted_16$1={class:"code-container"},_hoisted_17$1={class:"code-block"},_hoisted_18$1={class:"language-javascript"},_hoisted_19$1={class:"code-example"},_hoisted_20$1={class:"code-header"},_hoisted_21$1={class:"code-title"},_hoisted_22$1={class:"code-container"},_hoisted_23$1={class:"code-block"},_hoisted_24$1={class:"language-python"},_hoisted_25$1={class:"code-example"},_hoisted_26$1={class:"code-header"},_hoisted_27$1={class:"code-title"},_hoisted_28$1={class:"code-container"},_hoisted_29$1={class:"code-block"},_hoisted_30={class:"language-java"},_hoisted_31={class:"code-example"},_hoisted_32={class:"code-header"},_hoisted_33={class:"code-title"},_hoisted_34={class:"code-container"},_hoisted_35={class:"code-block"},_hoisted_36={class:"language-go"},_sfc_main$7=defineComponent({name:"ApiExamples",__name:"index",props:{codeExamples:{},copyStatus:{},copyCode:{type:Function},switchCodeTab:{type:Function}},setup(__props){const{$t:$t2}=useTexts(),props=__props,onCopyCode=(codeId,code)=>{props.copyCode(codeId,code)},onSwitchCodeTab=tabId=>{props.switchCodeTab(tabId)};return(_ctx,_cache)=>(openBlock(),createBlock(ApiSection,{id:"api-examples",title:unref($t2)("docs.examples.title"),"icon-class":"fas fa-code"},{default:withCtx(()=>[createVNode(CodeTabs,{tabs:unref(codeTabs),default:"curl","onUpdate:tab":onSwitchCodeTab},{curl:withCtx(()=>[createBaseVNode("div",_hoisted_1$7,[createBaseVNode("div",_hoisted_2$7,[createBaseVNode("span",_hoisted_3$7,toDisplayString(unref($t2)("docs.examples.singleUpload")),1),createBaseVNode("button",{class:normalizeClass(["copy-btn",{copied:__props.copyStatus["curl-single"]}]),onClick:_cache[0]||(_cache[0]=$event=>onCopyCode("curl-single",__props.codeExamples.curl.single))},[createBaseVNode("i",{class:normalizeClass(__props.copyStatus["curl-single"]?"fas fa-check":"fas fa-copy")},null,2),createBaseVNode("span",null,toDisplayString(__props.copyStatus["curl-single"]?unref($t2)("docs.examples.copied"):unref($t2)("docs.examples.copy")),1)],2)]),createBaseVNode("div",_hoisted_4$5,[createBaseVNode("pre",_hoisted_5$5,[createBaseVNode("code",_hoisted_6$5,toDisplayString(__props.codeExamples.curl.single),1)])])])]),javascript:withCtx(()=>[createBaseVNode("div",_hoisted_7$5,[createBaseVNode("div",_hoisted_8$2,[createBaseVNode("span",_hoisted_9$2,toDisplayString(unref($t2)("docs.examples.singleUpload")),1),createBaseVNode("button",{class:normalizeClass(["copy-btn",{copied:__props.copyStatus["js-single"]}]),onClick:_cache[1]||(_cache[1]=$event=>onCopyCode("js-single",__props.codeExamples.javascript.single))},[createBaseVNode("i",{class:normalizeClass(__props.copyStatus["js-single"]?"fas fa-check":"fas fa-copy")},null,2),createBaseVNode("span",null,toDisplayString(__props.copyStatus["js-single"]?unref($t2)("docs.examples.copied"):unref($t2)("docs.examples.copy")),1)],2)]),createBaseVNode("div",_hoisted_10$2,[createBaseVNode("pre",_hoisted_11$2,[createBaseVNode("code",_hoisted_12$2,toDisplayString(__props.codeExamples.javascript.single),1)])])])]),nodejs:withCtx(()=>[createBaseVNode("div",_hoisted_13$2,[createBaseVNode("div",_hoisted_14$2,[createBaseVNode("span",_hoisted_15$1,toDisplayString(unref($t2)("docs.examples.singleUpload")),1),createBaseVNode("button",{class:normalizeClass(["copy-btn",{copied:__props.copyStatus["nodejs-single"]}]),onClick:_cache[2]||(_cache[2]=$event=>onCopyCode("nodejs-single",__props.codeExamples.nodejs.single))},[createBaseVNode("i",{class:normalizeClass(__props.copyStatus["nodejs-single"]?"fas fa-check":"fas fa-copy")},null,2),createBaseVNode("span",null,toDisplayString(__props.copyStatus["nodejs-single"]?unref($t2)("docs.examples.copied"):unref($t2)("docs.examples.copy")),1)],2)]),createBaseVNode("div",_hoisted_16$1,[createBaseVNode("pre",_hoisted_17$1,[createBaseVNode("code",_hoisted_18$1,toDisplayString(__props.codeExamples.nodejs.single),1)])])])]),python:withCtx(()=>[createBaseVNode("div",_hoisted_19$1,[createBaseVNode("div",_hoisted_20$1,[createBaseVNode("span",_hoisted_21$1,toDisplayString(unref($t2)("docs.examples.basicUpload")),1),createBaseVNode("button",{class:normalizeClass(["copy-btn",{copied:__props.copyStatus["py-basic"]}]),onClick:_cache[3]||(_cache[3]=$event=>onCopyCode("py-basic",__props.codeExamples.python.basic))},[createBaseVNode("i",{class:normalizeClass(__props.copyStatus["py-basic"]?"fas fa-check":"fas fa-copy")},null,2),createBaseVNode("span",null,toDisplayString(__props.copyStatus["py-basic"]?unref($t2)("docs.examples.copied"):unref($t2)("docs.examples.copy")),1)],2)]),createBaseVNode("div",_hoisted_22$1,[createBaseVNode("pre",_hoisted_23$1,[createBaseVNode("code",_hoisted_24$1,toDisplayString(__props.codeExamples.python.basic),1)])])])]),java:withCtx(()=>[createBaseVNode("div",_hoisted_25$1,[createBaseVNode("div",_hoisted_26$1,[createBaseVNode("span",_hoisted_27$1,toDisplayString(unref($t2)("docs.examples.basicUpload")),1),createBaseVNode("button",{class:normalizeClass(["copy-btn",{copied:__props.copyStatus["java-basic"]}]),onClick:_cache[4]||(_cache[4]=$event=>onCopyCode("java-basic",__props.codeExamples.java.basic))},[createBaseVNode("i",{class:normalizeClass(__props.copyStatus["java-basic"]?"fas fa-check":"fas fa-copy")},null,2),createBaseVNode("span",null,toDisplayString(__props.copyStatus["java-basic"]?unref($t2)("docs.examples.copied"):unref($t2)("docs.examples.copy")),1)],2)]),createBaseVNode("div",_hoisted_28$1,[createBaseVNode("pre",_hoisted_29$1,[createBaseVNode("code",_hoisted_30,toDisplayString(__props.codeExamples.java.basic),1)])])])]),go:withCtx(()=>[createBaseVNode("div",_hoisted_31,[createBaseVNode("div",_hoisted_32,[createBaseVNode("span",_hoisted_33,toDisplayString(unref($t2)("docs.examples.singleUpload")),1),createBaseVNode("button",{class:normalizeClass(["copy-btn",{copied:__props.copyStatus["go-single"]}]),onClick:_cache[5]||(_cache[5]=$event=>onCopyCode("go-single",__props.codeExamples.go.single))},[createBaseVNode("i",{class:normalizeClass(__props.copyStatus["go-single"]?"fas fa-check":"fas fa-copy")},null,2),createBaseVNode("span",null,toDisplayString(__props.copyStatus["go-single"]?unref($t2)("docs.examples.copied"):unref($t2)("docs.examples.copy")),1)],2)]),createBaseVNode("div",_hoisted_34,[createBaseVNode("pre",_hoisted_35,[createBaseVNode("code",_hoisted_36,toDisplayString(__props.codeExamples.go.single),1)])])])]),_:1},8,["tabs"])]),_:1},8,["title"]))}});const ApiExamples=_export_sfc(_sfc_main$7,[["__scopeId","data-v-4e5403c3"]]),_sfc_main$6=defineComponent({name:"AuthenticationSection",components:{ApiSection},setup(){const{$t:$t2}=useTexts();return{$t:$t2}}});const _hoisted_1$6={class:"content-card"},_hoisted_2$6={class:"card-header"},_hoisted_3$6={class:"badge"},_hoisted_4$4={class:"examples-grid"},_hoisted_5$4={class:"example-item"},_hoisted_6$4={class:"example-item"},_hoisted_7$4={class:"tips-card"};function _sfc_render$3(_ctx,_cache,$props,$setup,$data,$options){const _component_ApiSection=resolveComponent("ApiSection");return openBlock(),createBlock(_component_ApiSection,{id:"api-authentication",title:_ctx.$t("docs.auth.title"),"icon-class":"fas fa-key"},{default:withCtx(()=>[createBaseVNode("div",_hoisted_1$6,[createBaseVNode("div",_hoisted_2$6,[_cache[0]||(_cache[0]=createBaseVNode("i",{class:"fas fa-key"},null,-1)),createBaseVNode("h4",null,toDisplayString(_ctx.$t("docs.auth.apiKeyAuth")),1),createBaseVNode("span",_hoisted_3$6,toDisplayString(_ctx.$t("docs.auth.recommended")),1)]),createBaseVNode("p",null,toDisplayString(_ctx.$t("docs.auth.twoMethods")),1),createBaseVNode("div",_hoisted_4$4,[createBaseVNode("div",_hoisted_5$4,[createBaseVNode("h5",null,toDisplayString(_ctx.$t("docs.auth.httpHeader")),1),_cache[1]||(_cache[1]=createBaseVNode("div",{class:"code-snippet"},"x-pixelpunk-key: YOUR_API_KEY",-1))]),createBaseVNode("div",_hoisted_6$4,[createBaseVNode("h5",null,toDisplayString(_ctx.$t("docs.auth.queryParam")),1),_cache[2]||(_cache[2]=createBaseVNode("div",{class:"code-snippet"},"?key=YOUR_API_KEY",-1))])])]),createBaseVNode("div",_hoisted_7$4,[createBaseVNode("h4",null,[_cache[3]||(_cache[3]=createBaseVNode("i",{class:"fas fa-shield-alt"},null,-1)),createTextVNode(toDisplayString(_ctx.$t("docs.auth.securityTips")),1)]),createBaseVNode("ul",null,[createBaseVNode("li",null,toDisplayString(_ctx.$t("docs.auth.tip1")),1),createBaseVNode("li",null,toDisplayString(_ctx.$t("docs.auth.tip2")),1),createBaseVNode("li",null,toDisplayString(_ctx.$t("docs.auth.tip3")),1),createBaseVNode("li",null,toDisplayString(_ctx.$t("docs.auth.tip4")),1)])])]),_:1},8,["title"])}const AuthenticationSection=_export_sfc(_sfc_main$6,[["render",_sfc_render$3],["__scopeId","data-v-2e5c34c0"]]),_hoisted_1$5={class:"code-example"},_hoisted_2$5={class:"code-header"},_hoisted_3$5={class:"code-title"},_hoisted_4$3=["title"],_hoisted_5$3={class:"copy-text"},_hoisted_6$3={class:"code-container"},_hoisted_7$3=["innerHTML"],_sfc_main$5=defineComponent({name:"CodeExample",__name:"index",props:{title:{},language:{default:"javascript"},marginBottom:{type:Boolean,default:!1}},setup(__props){const{$t:$t2}=useTexts(),props=__props,slots=useSlots(),codeElement=ref(),rawCode=ref(""),copied=ref(!1);onMounted(()=>{if(slots.default){const slotContent=slots.default();slotContent&&slotContent[0]&&slotContent[0].children&&(rawCode.value=slotContent[0].children.toString())}});const cleanCode=computed(()=>{if(!rawCode.value)return"";const lines=rawCode.value.split(`
`);for(;lines.length>0&&lines[0].trim()==="";)lines.shift();for(;lines.length>0&&lines[lines.length-1].trim()==="";)lines.pop();let code=lines.join(`
`);if(props.language==="json")try{const parsed=JSON.parse(code);code=JSON.stringify(parsed,null,2)}catch{}return code}),highlightedCode=computed(()=>{const code=cleanCode.value;if(!code)return"";const escaped=code.replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/>/g,"&gt;");switch(props.language){case"json":return highlightJson(escaped);case"javascript":return highlightJavaScript(escaped);case"bash":return highlightBash(escaped);case"python":return highlightPython(escaped);case"go":return highlightGo(escaped);default:return escaped}}),highlightJson=code=>code.replace(/"([^"]+)":/g,'<span class="token-key">"$1":</span>').replace(/:\s*"([^"]*)"/g,': <span class="token-string">"$1"</span>').replace(/:\s*(\d+)/g,': <span class="token-number">$1</span>').replace(/:\s*(true|false|null)/g,': <span class="token-boolean">$1</span>'),highlightJavaScript=code=>code.replace(/\b(async|function|const|let|var|if|else|for|while|return|try|catch|throw|import|export|from|class|new|this|super|await)\b/g,'<span class="token-keyword">$1</span>').replace(/'([^']*)'/g,`<span class="token-string">'$1'</span>`).replace(/"([^"]*)"/g,'<span class="token-string">"$1"</span>').replace(/`([^`]*)`/g,'<span class="token-template">`$1`</span>').replace(/\/\/.*$/gm,'<span class="token-comment">$&</span>').replace(/\b(\d+)\b/g,'<span class="token-number">$1</span>').replace(/\b(console|document|window|navigator)\b/g,'<span class="token-builtin">$1</span>'),highlightBash=code=>code.replace(/^#.*$/gm,'<span class="token-comment">$&</span>').replace(/\b(curl|npm|git|cd|ls|mkdir|rm|cp|mv|cat|grep|find|chmod|sudo|wget|tar|zip|unzip)\b/g,'<span class="token-command">$1</span>').replace(/-{1,2}[a-zA-Z-]+/g,'<span class="token-flag">$&</span>').replace(/"([^"]*)"/g,'<span class="token-string">"$1"</span>').replace(/'([^']*)'/g,`<span class="token-string">'$1'</span>`),highlightPython=code=>code.replace(/\b(def|class|import|from|if|elif|else|for|while|try|except|finally|with|as|return|yield|break|continue|pass|lambda|and|or|not|in|is|None|True|False|self)\b/g,'<span class="token-keyword">$1</span>').replace(/#.*$/gm,'<span class="token-comment">$&</span>').replace(/"""[\s\S]*?"""/g,'<span class="token-docstring">$&</span>').replace(/"([^"]*)"/g,'<span class="token-string">"$1"</span>').replace(/'([^']*)'/g,`<span class="token-string">'$1'</span>`).replace(/\b(\d+\.?\d*)\b/g,'<span class="token-number">$1</span>').replace(/\b(print|len|range|enumerate|zip|map|filter|str|int|float|list|dict|set|tuple)\b/g,'<span class="token-builtin">$1</span>'),highlightGo=code=>code.replace(/\b(package|import|func|var|const|type|struct|interface|if|else|for|range|switch|case|default|return|defer|go|chan|select|break|continue|fallthrough|goto)\b/g,'<span class="token-keyword">$1</span>').replace(/\/\/.*$/gm,'<span class="token-comment">$&</span>').replace(/\/\*[\s\S]*?\*\//g,'<span class="token-comment">$&</span>').replace(/"([^"]*)"/g,'<span class="token-string">"$1"</span>').replace(/`([^`]*)`/g,'<span class="token-string">`$1`</span>').replace(/\b(\d+)\b/g,'<span class="token-number">$1</span>').replace(/\b(string|int|int64|float64|bool|byte|rune|error|fmt|http|json|time|os|io)\b/g,'<span class="token-type">$1</span>'),copyCode=()=>{const code=cleanCode.value;code&&navigator.clipboard.writeText(code).then(()=>{copied.value=!0,setTimeout(()=>{copied.value=!1},2e3)}).catch(err=>{console.error("Copy failed:",err)})};return(_ctx,_cache)=>(openBlock(),createElementBlock("div",_hoisted_1$5,[createBaseVNode("div",_hoisted_2$5,[createBaseVNode("span",_hoisted_3$5,toDisplayString(__props.title),1),createBaseVNode("button",{class:normalizeClass(["copy-btn",{copied:copied.value}]),title:unref($t2)("docs.examples.copy"),onClick:copyCode},[createBaseVNode("i",{class:normalizeClass(copied.value?"fas fa-check":"fas fa-copy")},null,2),createBaseVNode("span",_hoisted_5$3,toDisplayString(copied.value?unref($t2)("docs.examples.copied"):unref($t2)("docs.examples.copy")),1)],10,_hoisted_4$3)]),createBaseVNode("div",_hoisted_6$3,[createBaseVNode("pre",{class:normalizeClass(["code-block",{"mb-4":__props.marginBottom}])},[createBaseVNode("code",{ref_key:"codeElement",ref:codeElement,class:normalizeClass(`language-${__props.language}`),innerHTML:highlightedCode.value},null,10,_hoisted_7$3)],2)])]))}});const CodeExample=_export_sfc(_sfc_main$5,[["__scopeId","data-v-dce694bb"]]),_sfc_main$4=defineComponent({name:"UploadApiSection",components:{ApiSection,CodeExample},props:{supportedFormats:{type:Array,required:!0},jsonResponseExamples:{type:Object,required:!0}},setup(){const{$t:$t2}=useTexts();return{$t:$t2}}});const _hoisted_1$4={class:"api-endpoint"},_hoisted_2$4={class:"endpoint-content"},_hoisted_3$4={class:"params-table"},_hoisted_4$2={class:"param-row header"},_hoisted_5$2={class:"param-row"},_hoisted_6$2={class:"required"},_hoisted_7$2={class:"param-row"},_hoisted_8$1={class:"required"},_hoisted_9$1={class:"param-row"},_hoisted_10$1={class:"param-row"},_hoisted_11$1={class:"param-row"},_hoisted_12$1={class:"param-row"},_hoisted_13$1={class:"param-note"},_hoisted_14$1={class:"format-tags"};function _sfc_render$2(_ctx,_cache,$props,$setup,$data,$options){const _component_CodeExample=resolveComponent("CodeExample"),_component_ApiSection=resolveComponent("ApiSection");return openBlock(),createBlock(_component_ApiSection,{id:"api-upload",title:_ctx.$t("docs.upload.title"),"icon-class":"fas fa-upload"},{default:withCtx(()=>[createBaseVNode("div",_hoisted_1$4,[_cache[12]||(_cache[12]=createBaseVNode("div",{class:"endpoint-header"},[createBaseVNode("span",{class:"method post"},"POST"),createBaseVNode("code",null,"/api/v1/external/upload")],-1)),createBaseVNode("div",_hoisted_2$4,[createBaseVNode("h4",null,toDisplayString(_ctx.$t("docs.upload.requestParams")),1),createBaseVNode("div",_hoisted_3$4,[createBaseVNode("div",_hoisted_4$2,[createBaseVNode("span",null,toDisplayString(_ctx.$t("docs.upload.paramName")),1),createBaseVNode("span",null,toDisplayString(_ctx.$t("docs.upload.paramType")),1),createBaseVNode("span",null,toDisplayString(_ctx.$t("docs.upload.paramRequired")),1),createBaseVNode("span",null,toDisplayString(_ctx.$t("docs.upload.paramDesc")),1)]),createBaseVNode("div",_hoisted_5$2,[_cache[0]||(_cache[0]=createBaseVNode("code",null,"file",-1)),_cache[1]||(_cache[1]=createBaseVNode("span",null,"File",-1)),createBaseVNode("span",_hoisted_6$2,toDisplayString(_ctx.$t("docs.upload.yes"))+"*",1),createBaseVNode("span",null,toDisplayString(_ctx.$t("docs.upload.fileDesc")),1)]),createBaseVNode("div",_hoisted_7$2,[_cache[2]||(_cache[2]=createBaseVNode("code",null,"files[]",-1)),_cache[3]||(_cache[3]=createBaseVNode("span",null,"File[]",-1)),createBaseVNode("span",_hoisted_8$1,toDisplayString(_ctx.$t("docs.upload.yes"))+"*",1),createBaseVNode("span",null,toDisplayString(_ctx.$t("docs.upload.filesDesc")),1)]),createBaseVNode("div",_hoisted_9$1,[_cache[4]||(_cache[4]=createBaseVNode("code",null,"folderId",-1)),_cache[5]||(_cache[5]=createBaseVNode("span",null,"String",-1)),createBaseVNode("span",null,toDisplayString(_ctx.$t("docs.upload.no")),1),createBaseVNode("span",null,toDisplayString(_ctx.$t("docs.upload.folderIdDesc")),1)]),createBaseVNode("div",_hoisted_10$1,[_cache[6]||(_cache[6]=createBaseVNode("code",null,"filePath",-1)),_cache[7]||(_cache[7]=createBaseVNode("span",null,"String",-1)),createBaseVNode("span",null,toDisplayString(_ctx.$t("docs.upload.no")),1),createBaseVNode("span",null,toDisplayString(_ctx.$t("docs.upload.filePathDesc")),1)]),createBaseVNode("div",_hoisted_11$1,[_cache[8]||(_cache[8]=createBaseVNode("code",null,"access_level",-1)),_cache[9]||(_cache[9]=createBaseVNode("span",null,"String",-1)),createBaseVNode("span",null,toDisplayString(_ctx.$t("docs.upload.no")),1),createBaseVNode("span",null,toDisplayString(_ctx.$t("docs.upload.accessLevelDesc")),1)]),createBaseVNode("div",_hoisted_12$1,[_cache[10]||(_cache[10]=createBaseVNode("code",null,"optimize",-1)),_cache[11]||(_cache[11]=createBaseVNode("span",null,"String",-1)),createBaseVNode("span",null,toDisplayString(_ctx.$t("docs.upload.no")),1),createBaseVNode("span",null,toDisplayString(_ctx.$t("docs.upload.optimizeDesc")),1)])]),createBaseVNode("p",_hoisted_13$1,toDisplayString(_ctx.$t("docs.upload.paramNote")),1),createBaseVNode("h4",null,toDisplayString(_ctx.$t("docs.upload.supportedFormats")),1),createBaseVNode("div",_hoisted_14$1,[(openBlock(!0),createElementBlock(Fragment,null,renderList(_ctx.supportedFormats,format=>(openBlock(),createElementBlock("span",{key:format,class:"format-tag"},toDisplayString(format),1))),128))]),createBaseVNode("h4",null,toDisplayString(_ctx.$t("docs.upload.successResponse")),1),createVNode(_component_CodeExample,{title:_ctx.$t("docs.upload.singleSuccess"),language:"json"},{default:withCtx(()=>[createTextVNode(toDisplayString(_ctx.jsonResponseExamples.single),1)]),_:1},8,["title"]),createVNode(_component_CodeExample,{title:_ctx.$t("docs.upload.batchSuccess"),language:"json"},{default:withCtx(()=>[createTextVNode(toDisplayString(_ctx.jsonResponseExamples.batch),1)]),_:1},8,["title"])])])]),_:1},8,["title"])}const UploadApiSection=_export_sfc(_sfc_main$4,[["render",_sfc_render$2],["__scopeId","data-v-4d2830a0"]]),_sfc_main$3=defineComponent({name:"ApiLimitsSection",components:{ApiSection},props:{supportedFormats:{type:Array,required:!0}},setup(props){const{$t:$t2}=useTexts(),fileFormatText=computed(()=>{const formatsStr=Array.isArray(props.supportedFormats)?props.supportedFormats.join(", "):"";return $t2("docs.limits.fileFormat",{formats:formatsStr})});return{$t:$t2,fileFormatText}}});const _hoisted_1$3={class:"limits-grid"},_hoisted_2$3={class:"limit-card"},_hoisted_3$3={class:"limit-card"},_hoisted_4$1={class:"error-codes"},_hoisted_5$1={class:"error-list"},_hoisted_6$1={class:"error-item"},_hoisted_7$1={class:"error-details"},_hoisted_8={class:"error-title"},_hoisted_9={class:"error-solution"},_hoisted_10={class:"error-item"},_hoisted_11={class:"error-details"},_hoisted_12={class:"error-title"},_hoisted_13={class:"error-solution"},_hoisted_14={class:"error-item"},_hoisted_15={class:"error-details"},_hoisted_16={class:"error-title"},_hoisted_17={class:"error-solution"},_hoisted_18={class:"error-item"},_hoisted_19={class:"error-details"},_hoisted_20={class:"error-title"},_hoisted_21={class:"error-solution"},_hoisted_22={class:"error-item"},_hoisted_23={class:"error-details"},_hoisted_24={class:"error-title"},_hoisted_25={class:"error-solution"},_hoisted_26={class:"error-item"},_hoisted_27={class:"error-details"},_hoisted_28={class:"error-title"},_hoisted_29={class:"error-solution"};function _sfc_render$1(_ctx,_cache,$props,$setup,$data,$options){const _component_ApiSection=resolveComponent("ApiSection");return openBlock(),createBlock(_component_ApiSection,{id:"api-limits",title:_ctx.$t("docs.limits.title"),"icon-class":"fas fa-exclamation-circle"},{default:withCtx(()=>[createBaseVNode("div",_hoisted_1$3,[createBaseVNode("div",_hoisted_2$3,[createBaseVNode("h4",null,[_cache[0]||(_cache[0]=createBaseVNode("i",{class:"fas fa-file-image"},null,-1)),createTextVNode(toDisplayString(_ctx.$t("docs.limits.fileLimit")),1)]),createBaseVNode("ul",null,[createBaseVNode("li",null,toDisplayString(_ctx.fileFormatText),1),createBaseVNode("li",null,toDisplayString(_ctx.$t("docs.limits.singleFileSize")),1),createBaseVNode("li",null,toDisplayString(_ctx.$t("docs.limits.batchUpload")),1),createBaseVNode("li",null,toDisplayString(_ctx.$t("docs.limits.autoOptimize")),1)])]),createBaseVNode("div",_hoisted_3$3,[createBaseVNode("h4",null,[_cache[1]||(_cache[1]=createBaseVNode("i",{class:"fas fa-key"},null,-1)),createTextVNode(toDisplayString(_ctx.$t("docs.limits.apiKeyLimit")),1)]),createBaseVNode("ul",null,[createBaseVNode("li",null,toDisplayString(_ctx.$t("docs.limits.storageCapacity")),1),createBaseVNode("li",null,toDisplayString(_ctx.$t("docs.limits.uploadCount")),1),createBaseVNode("li",null,toDisplayString(_ctx.$t("docs.limits.singleFileSizeLimit")),1),createBaseVNode("li",null,toDisplayString(_ctx.$t("docs.limits.defaultFolder")),1)])])]),createBaseVNode("div",_hoisted_4$1,[createBaseVNode("h4",null,[_cache[2]||(_cache[2]=createBaseVNode("i",{class:"fas fa-exclamation-triangle"},null,-1)),createTextVNode(toDisplayString(_ctx.$t("docs.limits.commonErrors")),1)]),createBaseVNode("div",_hoisted_5$1,[createBaseVNode("div",_hoisted_6$1,[_cache[3]||(_cache[3]=createBaseVNode("span",{class:"error-code"},"102",-1)),createBaseVNode("div",_hoisted_7$1,[createBaseVNode("span",_hoisted_8,toDisplayString(_ctx.$t("docs.limits.error102Title")),1),createBaseVNode("span",_hoisted_9,toDisplayString(_ctx.$t("docs.limits.error102Solution")),1)])]),createBaseVNode("div",_hoisted_10,[_cache[4]||(_cache[4]=createBaseVNode("span",{class:"error-code"},"4000",-1)),createBaseVNode("div",_hoisted_11,[createBaseVNode("span",_hoisted_12,toDisplayString(_ctx.$t("docs.limits.error4000Title")),1),createBaseVNode("span",_hoisted_13,toDisplayString(_ctx.$t("docs.limits.error4000Solution")),1)])]),createBaseVNode("div",_hoisted_14,[_cache[5]||(_cache[5]=createBaseVNode("span",{class:"error-code"},"4001",-1)),createBaseVNode("div",_hoisted_15,[createBaseVNode("span",_hoisted_16,toDisplayString(_ctx.$t("docs.limits.error4001Title")),1),createBaseVNode("span",_hoisted_17,toDisplayString(_ctx.$t("docs.limits.error4001Solution")),1)])]),createBaseVNode("div",_hoisted_18,[_cache[6]||(_cache[6]=createBaseVNode("span",{class:"error-code"},"4008",-1)),createBaseVNode("div",_hoisted_19,[createBaseVNode("span",_hoisted_20,toDisplayString(_ctx.$t("docs.limits.error4008Title")),1),createBaseVNode("span",_hoisted_21,toDisplayString(_ctx.$t("docs.limits.error4008Solution")),1)])]),createBaseVNode("div",_hoisted_22,[_cache[7]||(_cache[7]=createBaseVNode("span",{class:"error-code"},"4009",-1)),createBaseVNode("div",_hoisted_23,[createBaseVNode("span",_hoisted_24,toDisplayString(_ctx.$t("docs.limits.error4009Title")),1),createBaseVNode("span",_hoisted_25,toDisplayString(_ctx.$t("docs.limits.error4009Solution")),1)])]),createBaseVNode("div",_hoisted_26,[_cache[8]||(_cache[8]=createBaseVNode("span",{class:"error-code"},"4010",-1)),createBaseVNode("div",_hoisted_27,[createBaseVNode("span",_hoisted_28,toDisplayString(_ctx.$t("docs.limits.error4010Title")),1),createBaseVNode("span",_hoisted_29,toDisplayString(_ctx.$t("docs.limits.error4010Solution")),1)])])])])]),_:1},8,["title"])}const ApiLimitsSection=_export_sfc(_sfc_main$3,[["render",_sfc_render$1],["__scopeId","data-v-7bb06df3"]]),_sfc_main$2=defineComponent({name:"FaqSection",components:{ApiSection},setup(){const{$t:$t2}=useTexts(),openFaqIndex=ref(null),faqList=$t2("docs.faq.items");return{openFaqIndex,faqList,toggleFaq:index2=>{openFaqIndex.value=openFaqIndex.value===index2?null:index2}}}});const _hoisted_1$2={class:"faq-accordion"},_hoisted_2$2=["onClick"],_hoisted_3$2={class:"faq-title-wrapper"},_hoisted_4={class:"faq-title"},_hoisted_5={class:"faq-toggle"},_hoisted_6={class:"faq-content"},_hoisted_7={class:"faq-answer"};function _sfc_render(_ctx,_cache,$props,$setup,$data,$options){const _component_ApiSection=resolveComponent("ApiSection");return openBlock(),createBlock(_component_ApiSection,{id:"api-faq",title:_ctx.$t("docs.faq.title"),"icon-class":"fas fa-question-circle"},{default:withCtx(()=>[createBaseVNode("div",_hoisted_1$2,[(openBlock(!0),createElementBlock(Fragment,null,renderList(_ctx.faqList,(faq,index2)=>(openBlock(),createElementBlock("div",{key:index2,class:normalizeClass(["faq-accordion-item",{active:_ctx.openFaqIndex===index2}])},[createBaseVNode("div",{class:"faq-header",onClick:$event=>_ctx.toggleFaq(index2)},[createBaseVNode("div",_hoisted_3$2,[createBaseVNode("i",{class:normalizeClass([faq.icon,"faq-icon"])},null,2),createBaseVNode("h4",_hoisted_4,toDisplayString(faq.question),1)]),createBaseVNode("div",_hoisted_5,[createBaseVNode("i",{class:normalizeClass(["fas fa-chevron-down",{rotated:_ctx.openFaqIndex===index2}])},null,2)])],8,_hoisted_2$2),createVNode(Transition,{name:"faq-content"},{default:withCtx(()=>[withDirectives(createBaseVNode("div",_hoisted_6,[createBaseVNode("div",_hoisted_7,[createBaseVNode("p",null,toDisplayString(faq.answer),1)])],512),[[vShow,_ctx.openFaqIndex===index2]])]),_:2},1024)],2))),128))])]),_:1},8,["title"])}const FaqSection=_export_sfc(_sfc_main$2,[["render",_sfc_render],["__scopeId","data-v-a28382c4"]]),_hoisted_1$1={class:"tester-features"},_hoisted_2$1={class:"feature-icon"},_hoisted_3$1={class:"feature-content"},_sfc_main$1=defineComponent({name:"ApiTesterFeatures",__name:"index",props:{features:{default:()=>[]}},setup(__props){return(_ctx,_cache)=>(openBlock(),createElementBlock("div",_hoisted_1$1,[(openBlock(!0),createElementBlock(Fragment,null,renderList(__props.features,feature=>(openBlock(),createElementBlock("div",{key:feature.id,class:"feature-item"},[createBaseVNode("div",_hoisted_2$1,[createBaseVNode("i",{class:normalizeClass(feature.icon)},null,2)]),createBaseVNode("div",_hoisted_3$1,[createBaseVNode("h5",null,toDisplayString(feature.title),1),createBaseVNode("p",null,toDisplayString(feature.description),1)])]))),128))]))}});const ApiTesterFeatures=_export_sfc(_sfc_main$1,[["__scopeId","data-v-72538e96"]]);function useCodeExamples(currentDomain){const codeExamples=computed(()=>({curl:{single:`# 基础上传
curl -X POST \\
  "${currentDomain.value}/api/v1/external/upload" \\
  -H "x-pixelpunk-key: YOUR_API_KEY" \\
  -F "file=@image.jpg"

# 带参数上传
curl -X POST \\
  "${currentDomain.value}/api/v1/external/upload" \\
  -H "x-pixelpunk-key: YOUR_API_KEY" \\
  -F "file=@image.jpg" \\
  -F "access_level=public" \\
  -F "optimize=true" \\
  -F "filePath=projects/website"

# 指定文件夹ID上传
curl -X POST \\
  "${currentDomain.value}/api/v1/external/upload" \\
  -H "x-pixelpunk-key: YOUR_API_KEY" \\
  -F "file=@image.jpg" \\
  -F "folderId=folder_12345" \\
  -F "optimize=1"`,batch:`# 批量上传多个文件
curl -X POST \\
  "${currentDomain.value}/api/v1/external/upload" \\
  -H "x-pixelpunk-key: YOUR_API_KEY" \\
  -F "files[]=@image1.jpg" \\
  -F "files[]=@image2.png" \\
  -F "files[]=@image3.gif" \\
  -F "filePath=batch/$(date +%Y%m%d)"

# 支持的多种字段名格式
curl -X POST \\
  "${currentDomain.value}/api/v1/external/upload" \\
  -H "x-pixelpunk-key: YOUR_API_KEY" \\
  -F "images[]=@photo1.jpg" \\
  -F "images[]=@photo2.jpg" \\
  -F "access_level=public"`},javascript:{single:`// 基础上传函数
async function uploadFile(file, options = {}) {
  const formData = new FormData();
  formData.append('file', file);

  if (options.access_level) formData.append('access_level', options.access_level);
  if (options.optimize) formData.append('optimize', options.optimize ? 'true' : 'false');
  if (options.filePath) formData.append('filePath', options.filePath);
  if (options.folderId) formData.append('folderId', options.folderId);

  try {
    const response = await fetch('${currentDomain.value}/api/v1/external/upload', {
      method: 'POST',
      headers: {
        'x-pixelpunk-key': 'YOUR_API_KEY'
      },
      body: formData
    });

    const result = await response.json();

    if (result.code === 0) {
      if (result.data.oversized_files) {
      }
      if (result.data.upload_errors) {
      }

      return result.data;
    } else {
      throw new Error(result.message);
    }
  } catch (error) {
    console.error('Upload failed:', error);
    throw error;
  }
}

async function uploadMultipleImages(files, options = {}) {
  const formData = new FormData();

  const fieldName = options.fieldName || 'files[]';
  files.forEach(file => {
    formData.append(fieldName, file);
  });

  if (options.filePath) formData.append('filePath', options.filePath);
  if (options.access_level) formData.append('access_level', options.access_level);
  if (options.optimize) formData.append('optimize', options.optimize ? 'true' : 'false');

}`},nodejs:{single:`// Node.js + Express + multer 示例
const express = require('express');
const multer = require('multer');
const FormData = require('form-data');
const fs = require('fs');
const axios = require('axios');

const app = express();

const upload = multer({ dest: 'uploads/' });

async function uploadToAPI(filePath, options = {}) {
  try {
    const formData = new FormData();
    formData.append('file', fs.createReadStream(filePath));

    if (options.access_level) formData.append('access_level', options.access_level);
    if (options.optimize) formData.append('optimize', options.optimize.toString());
    if (options.filePath) formData.append('filePath', options.filePath);
    if (options.folderId) formData.append('folderId', options.folderId);

    const response = await axios.post('${currentDomain.value}/api/v1/external/upload', formData, {
      headers: {
        'x-pixelpunk-key': 'YOUR_API_KEY',
        ...formData.getHeaders()
      },
      timeout: TIMING.REQUEST.LONG_TIMEOUT
    });

    return response.data;
  } catch (error) {
    console.error('Upload failed:', error.response?.data || error.message);
    throw error;
  }
}

async function uploadMultipleToAPI(filePaths, options = {}) {
  try {
    const formData = new FormData();

    filePaths.forEach(filePath => {
      formData.append('files[]', fs.createReadStream(filePath));
    });

    if (options.access_level) formData.append('access_level', options.access_level);
    if (options.optimize) formData.append('optimize', options.optimize.toString());
    if (options.filePath) formData.append('filePath', options.filePath);
    if (options.folderId) formData.append('folderId', options.folderId);

    const response = await axios.post('${currentDomain.value}/api/v1/external/upload', formData, {
      headers: {
        'x-pixelpunk-key': 'YOUR_API_KEY',
        ...formData.getHeaders()
      },
      timeout: 60000
    });

    return response.data;
  } catch (error) {
    console.error('Batch upload failed:', error.response?.data || error.message);
    throw error;
  }
}

app.post('/upload', upload.single('file'), async (req, res) => {
  try {
    const result = await uploadToAPI(req.file.path, {
      access_level: req.body.access_level || 'private',
      optimize: req.body.optimize === 'true',
      filePath: req.body.filePath,
      folderId: req.body.folderId
    });

    fs.unlinkSync(req.file.path);

    res.json(result);
  } catch (error) {
    res.status(500).json({ _error: error.message });
  }
});

app.post('/upload/batch', upload.array('files', 10), async (req, res) => {
  try {
    const filePaths = req.files.map(file => file.path);

    const result = await uploadMultipleToAPI(filePaths, {
      access_level: req.body.access_level || 'private',
      optimize: req.body.optimize === 'true',
      filePath: req.body.filePath,
      folderId: req.body.folderId
    });

    req.files.forEach(file => fs.unlinkSync(file.path));

    res.json(result);
  } catch (error) {
    res.status(500).json({ _error: error.message });
  }
});

app.listen(3000, () => {
});

`},python:{basic:`import requests

def upload_image(file_path, api_key):
    url = "${currentDomain.value}/api/v1/external/upload"

    headers = {
        "x-pixelpunk-key": api_key
    }

    files = {
        "file": open(file_path, "rb")
    }

    data = {
        "access_level": "private"
    }

    try:
        response = requests.post(url, headers=headers, files=files, data=data)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"上传失败: {e}")
        return None
    finally:
        files["file"].close()`},java:{basic:`import java.io.File;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;

public class ImageUploader {
    private static final String API_KEY = "YOUR_API_KEY";
    private static final String BASE_URL = "${currentDomain.value}";
    private static final String UPLOAD_URL = BASE_URL + "/api/v1/external/upload";

    public static void main(String[] args) {
        try {
            List<File> images = new ArrayList<>();
            images.add(new File("image1.jpg"));
            images.add(new File("image2.png"));
            images.add(new File("image3.gif"));

            List<String> imagePaths = images.stream()
                .map(File::getAbsolutePath)
                .collect(Collectors.toList());

            String folderId = "folder_12345";
            String filePath = "projects/website";
            String accessLevel = "private";
            boolean optimize = true;

            uploadSingleFile(images.get(0), accessLevel, optimize, filePath, folderId);

            uploadMultipleFiles(imagePaths, accessLevel, optimize, filePath, folderId);
        } catch (IOException | InterruptedException e) {
            e.printStackTrace();
        }
    }

    private static void uploadSingleFile(File file, String accessLevel, boolean optimize, String filePath, String folderId) throws IOException, InterruptedException {
        String command = String.format("curl -X POST "%s/api/v1/external/upload" -H "x-pixelpunk-key: %s" -F "file=@%s" -F "access_level=%s" -F "optimize=%s" -F "filePath=%s" -F "folderId=%s"",
            BASE_URL, API_KEY, file.getAbsolutePath(), accessLevel, optimize, filePath, folderId);

        Process process = Runtime.getRuntime().exec(command);
        process.waitFor(30, TimeUnit.SECONDS);

    }

    private static void uploadMultipleFiles(List<String> imagePaths, String accessLevel, boolean optimize, String filePath, String folderId) throws IOException, InterruptedException {
        String command = String.format("curl -X POST "%s/api/v1/external/upload" -H "x-pixelpunk-key: %s" -F "files[]=@%s" -F "access_level=%s" -F "optimize=%s" -F "filePath=%s" -F "folderId=%s"",
            BASE_URL, API_KEY, String.join("" -F "files[]=@" ", imagePaths), accessLevel, optimize, filePath, folderId);

        Process process = Runtime.getRuntime().exec(command);
        process.waitFor(30, TimeUnit.SECONDS);

    }
}`},go:{single:`package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type UploadResponse struct {
	Code    int    \`json:"code"\`
	Message string \`json:"message"\`
	Data    struct {
		Uploaded []UploadedFile \`json:"uploaded"\`
		OversizedFiles []string \`json:"oversized_files,omitempty"\`
		SizeLimit     string   \`json:"size_limit,omitempty"\`
		UploadErrors  []string \`json:"upload_errors,omitempty"\`
	} \`json:"data"\`
}

type UploadedFile struct {
	ID           string \`json:"id"\`
	URL          string \`json:"url"\`
	ThumbURL     string \`json:"thumb_url"\`
	OriginalName string \`json:"original_name"\`
	Size         int64  \`json:"size"\`
	Width        int    \`json:"width"\`
	Height       int    \`json:"height"\`
	Format       string \`json:"format"\`
	AccessLevel  string \`json:"access_level"\`
	CreatedAt    string \`json:"created_at"\`
}

type UploadOptions struct {
	APIKey     string
	AccessLevel string // public/private/protected
	Optimize   bool
	FilePath   string
	FolderID   string
}

type ImageUploader struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewImageUploader(baseURL string) *ImageUploader {
	return &ImageUploader{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (u *ImageUploader) UploadSingleFile(filePath string, options UploadOptions) (*UploadResponse, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 创建multipart form
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 添加文件字段
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("创建文件字段失败: %v", err)
	}

	// 拷贝文件内容
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("拷贝文件内容失败: %v", err)
	}

	// 添加其他字段
	if options.AccessLevel != "" {
		writer.WriteField("access_level", options.AccessLevel)
	}
	if options.Optimize {
		writer.WriteField("optimize", "true")
	} else {
		writer.WriteField("optimize", "false")
	}
	if options.FilePath != "" {
		writer.WriteField("filePath", options.FilePath)
	}
	if options.FolderID != "" {
		writer.WriteField("folderId", options.FolderID)
	}

	// 关闭writer
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("关闭writer失败: %v", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", u.BaseURL+"/api/v1/external/upload", &requestBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-pixelpunk-key", options.APIKey)

	// 发送请求
	resp, err := u.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析响应
	var uploadResp UploadResponse
	err = json.Unmarshal(body, &uploadResp)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if uploadResp.Code != 0 {
		return &uploadResp, fmt.Errorf("上传失败: %s", uploadResp.Message)
	}

	return &uploadResp, nil
}

func (u *ImageUploader) UploadMultipleFiles(filePaths []string, options UploadOptions) (*UploadResponse, error) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 添加多个文件
	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("无法打开文件 %s: %v", filePath, err)
		}

		part, err := writer.CreateFormFile("files[]", filepath.Base(filePath))
		if err != nil {
			file.Close()
			return nil, fmt.Errorf("创建文件字段失败: %v", err)
		}

		_, err = io.Copy(part, file)
		file.Close()
		if err != nil {
			return nil, fmt.Errorf("拷贝文件内容失败: %v", err)
		}
	}

	// 添加其他字段
	if options.AccessLevel != "" {
		writer.WriteField("access_level", options.AccessLevel)
	}
	if options.Optimize {
		writer.WriteField("optimize", "true")
	} else {
		writer.WriteField("optimize", "false")
	}
	if options.FilePath != "" {
		writer.WriteField("filePath", options.FilePath)
	}
	if options.FolderID != "" {
		writer.WriteField("folderId", options.FolderID)
	}

	err := writer.Close()
	if err != nil {
		return nil, fmt.Errorf("关闭writer失败: %v", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", u.BaseURL+"/api/v1/external/upload", &requestBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-pixelpunk-key", options.APIKey)

	// 发送请求
	resp, err := u.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析响应
	var uploadResp UploadResponse
	err = json.Unmarshal(body, &uploadResp)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &uploadResp, nil
}

func main() {
	// 创建上传客户端
	uploader := NewImageUploader("${currentDomain.value}")

	// 配置上传选项
	options := UploadOptions{
		APIKey:     "YOUR_API_KEY",
		Visibility: "public",
		Optimize:   true,
		FilePath:   "projects/website",
		FolderID:   "", // 可选，优先级高于FilePath
	}

	// 单文件上传示例
	fmt.Println("=== 单文件上传 ===")
	singleResult, err := uploader.UploadSingleFile("image.jpg", options)
	if err != nil {
		fmt.Printf("单文件上传失败: %v\\n", err)
	} else {
		fmt.Printf("单文件上传成功: %+v\\n", singleResult.Data.Uploaded)
	}

	// 批量上传示例
	fmt.Println("\\n=== 批量上传 ===")
	filePaths := []string{"image1.jpg", "image2.png", "image3.gif"}
	batchResult, err := uploader.UploadMultipleFiles(filePaths, options)
	if err != nil {
		fmt.Printf("批量上传失败: %v\\n", err)
	} else {
		fmt.Printf("批量上传结果: %d个文件成功\\n", len(batchResult.Data.Uploaded))
		
		// 显示成功上传的文件
		for i, uploaded := range batchResult.Data.Uploaded {
			fmt.Printf("[%d] %s -> %s\\n", i+1, uploaded.OriginalName, uploaded.URL)
		}
		
		// 显示警告信息
		if len(batchResult.Data.OversizedFiles) > 0 {
			fmt.Printf("超限文件: %v\\n", batchResult.Data.OversizedFiles)
		}
		if len(batchResult.Data.UploadErrors) > 0 {
			fmt.Printf("上传错误: %v\\n", batchResult.Data.UploadErrors)
		}
	}
}

func ConcurrentUpload(uploader *ImageUploader, filePaths []string, options UploadOptions) {
	type uploadResult struct {
		FilePath string
		Result   *UploadResponse
		Error    error
	}

	resultChan := make(chan uploadResult, len(filePaths))

	// 启动goroutine进行并发上传
	for _, filePath := range filePaths {
		go func(fp string) {
			result, err := uploader.UploadSingleFile(fp, options)
			resultChan <- uploadResult{
				FilePath: fp,
				Result:   result,
				Error:    err,
			}
		}(filePath)
	}

	// 收集结果
	for i := 0; i < len(filePaths); i++ {
		result := <-resultChan
		if result.Error != nil {
			fmt.Printf("文件 %s 上传失败: %v\\n", result.FilePath, result.Error)
		} else {
			fmt.Printf("文件 %s 上传成功: %s\\n",
				result.FilePath,
				result.Result.Data.Uploaded[0].URL)
		}
	}
}`}})),jsonResponseExamples=computed(()=>({single:`{
  "code": 0,
  "message": "上传成功",
  "data": {
    "uploaded": {
      "id": "img_1a2b3c4d5e",
      "url": "${currentDomain.value}/file/user_1/image.jpg",
      "thumb_url": "${currentDomain.value}/thumb/user_1/image.jpg",
      "original_name": "image.jpg",
      "size": 1024000,
      "width": 1920,
      "height": 1080,
      "format": "jpg",
      "access_level": "private",
      "created_at": "2024-01-15T08:30:45Z"
    }
  }
}`,batch:`{
  "code": 0,
  "message": "部分文件上传成功",
  "data": {
    "uploaded": [
      {
        "id": "img_1a2b3c4d5e",
        "url": "${currentDomain.value}/file/user_1/image1.jpg",
        "thumb_url": "${currentDomain.value}/thumb/user_1/image1.jpg",
        "original_name": "image1.jpg",
        "size": 1024000,
        "width": 1920,
        "height": 1080,
        "format": "jpg",
        "access_level": "private",
        "created_at": "2024-01-15T08:30:45Z"
      }
    ],
    "oversized_files": ["large_image.jpg"],
    "size_limit": "20.0MB",
    "upload_errors": ["corrupted.jpg: 文件损坏无法处理"]
  }
}`}));return{codeExamples,jsonResponseExamples}}const _hoisted_1={class:"docs-page"},_hoisted_2={class:"docs-content"},_hoisted_3={class:"section-intro"},_sfc_main=defineComponent({name:"DocsPage",__name:"index",setup(__props){const{activeSection,copyStatus,mobileMenuOpen,showBackToTop,readingProgress,currentDomain,toggleMobileMenu,handleMobileNavClick,scrollToSection,backToTop,copyCode,switchCodeTab}=useDocsPage(),{codeExamples,jsonResponseExamples}=useCodeExamples(currentDomain),apiTesterFeatures=getApiTesterFeatures();return(_ctx,_cache)=>{const _component_CyberParticleBackground=resolveComponent("CyberParticleBackground");return openBlock(),createElementBlock("div",_hoisted_1,[createVNode(_component_CyberParticleBackground,{theme:"docs","max-particles":60}),createVNode(DocsSidebar,{class:"desktop-only","nav-items":unref(navItems),"active-section":unref(activeSection),"reading-progress":unref(readingProgress),"scroll-to-section":unref(scrollToSection)},null,8,["nav-items","active-section","reading-progress","scroll-to-section"]),createBaseVNode("main",_hoisted_2,[createVNode(MobileHeader,{"nav-items":unref(navItems),"active-section":unref(activeSection),"mobile-menu-open":unref(mobileMenuOpen),"toggle-mobile-menu":unref(toggleMobileMenu),"handle-mobile-nav-click":unref(handleMobileNavClick)},null,8,["nav-items","active-section","mobile-menu-open","toggle-mobile-menu","handle-mobile-nav-click"]),createVNode(DocsIntro,{"scroll-to-section":unref(scrollToSection)},null,8,["scroll-to-section"]),createVNode(ApiOverview,{"current-domain":unref(currentDomain)},null,8,["current-domain"]),createVNode(AuthenticationSection),createVNode(UploadApiSection,{"supported-formats":unref(supportedFormats),"json-response-examples":unref(jsonResponseExamples)},null,8,["supported-formats","json-response-examples"]),createVNode(ApiLimitsSection,{"supported-formats":unref(supportedFormats)},null,8,["supported-formats"]),createVNode(ApiExamples,{"code-examples":unref(codeExamples),"copy-status":unref(copyStatus),"copy-code":unref(copyCode),"switch-code-tab":unref(switchCodeTab)},null,8,["code-examples","copy-status","copy-code","switch-code-tab"]),createVNode(ApiSection,{id:"api-tester",title:_ctx.$t("docs.tester.title"),"icon-class":"fas fa-flask"},{default:withCtx(()=>[createBaseVNode("div",_hoisted_3,[createBaseVNode("p",null,toDisplayString(_ctx.$t("docs.tester.description")),1)]),createVNode(ApiTesterFeatures,{features:unref(apiTesterFeatures)},null,8,["features"]),createVNode(ApiTester)]),_:1},8,["title"]),createVNode(FaqSection)]),createVNode(BackToTop,{"show-back-to-top":unref(showBackToTop),"back-to-top":unref(backToTop)},null,8,["show-back-to-top","back-to-top"])])}}});const index=_export_sfc(_sfc_main,[["__scopeId","data-v-6fcaf8c9"]]);export{index as default};
