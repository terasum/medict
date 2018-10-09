<template lang="html">
    <div class="me-setting">
      <Form :model="formdata" :label-width="80">
        <FormItem label="MDX文件">
            <Input v-model="formdata.mdx" @on-search="selectMDX" search enter-button="选择MDX" placeholder="Select mdx..." />
        </FormItem>
        <FormItem label="MDD文件">
            <Input v-model="formdata.mdd" @on-search="selectMDD" search enter-button="选择MDD" placeholder="Select mdd ..." />
        </FormItem>
        <FormItem label="附加JS文件">
            <Input v-model="formdata.js" @on-search="selectJS" search enter-button="选择JS" placeholder="Select js ..." />
        </FormItem>
        <FormItem label="附加CSS文件">
            <Input v-model="formdata.css" @on-search="selectCSS" search enter-button="选择CSS" placeholder="Select css ..." />
        </FormItem>
        <Button @click="save" type="primary"> 保存 </Button>
        <Button @click="reset"> 重置 </Button>
    </Form>
    </div>
</template>

<style lang="scss" scoped>
  .me-setting{
    height: 100%;
    overflow-y: auto;
    width: 100%;
    overflow-x: hidden;
    padding: 10px 20px;
  }
</style>

<script>
import { remote, ipcRenderer } from 'electron'
import fs from 'fs'
const dialog = remote.dialog

export default {
  data () {
    return {
      formdata: {
        mdx: this.$store.state.Query.mdx,
        mdd: this.$store.state.Query.mdd,
        js: this.$store.state.Query.js,
        css: this.$store.state.Query.css
      }
    }
  },
  methods: {
    selectMDX () {
      console.log('select mdx')
      const mdxFilePath = dialog.showOpenDialog({
        properties: ['openFile'],
        filters: [
          {name: 'mdict', extensions: ['mdx']},
          {name: 'All Files', extensions: ['*']}
        ]
      })
      console.log(mdxFilePath)
      if (mdxFilePath && mdxFilePath[0]) {
        this.formdata.mdx = mdxFilePath[0]
      }
    },
    selectMDD () {
      console.log('select mdd')
      const mddFilePath = dialog.showOpenDialog({
        properties: ['openFile'],
        filters: [
          {name: 'mdict', extensions: ['mdd']},
          {name: 'All Files', extensions: ['*']}
        ]
      })
      console.log(mddFilePath)
      if (mddFilePath && mddFilePath[0]) {
        this.formdata.mdd = mddFilePath[0]
      } else {
        this.formdata.mdd = ''
      }
    },
    selectJS () {
      console.log('select js')
      const jsFilePath = dialog.showOpenDialog({
        properties: ['openFile'],
        filters: [
          {name: 'mdict', extensions: ['js']},
          {name: 'All Files', extensions: ['*']}
        ]
      })
      console.log(jsFilePath)
      if (jsFilePath && jsFilePath[0]) {
        this.formdata.js = jsFilePath[0]
      } else {
        this.formdata.js = ''
      }
    },
    selectCSS () {
      console.log('select css')
      const cssFilePath = dialog.showOpenDialog({
        properties: ['openFile'],
        filters: [
          {name: 'mdict', extensions: ['css']},
          {name: 'All Files', extensions: ['*']}
        ]
      })
      console.log(cssFilePath)
      if (cssFilePath && cssFilePath[0]) {
        this.formdata.css = cssFilePath[0]
      } else {
        this.formdata.css = ''
      }
    },
    save () {
      console.log('save')
      console.log(this.formdata.mdd)
      console.log(this.formdata.mdd.endsWith('.mdd'))
      console.log(this.formdata.mdx)
      console.log(this.formdata.mdd.endsWith('.mdx'))
      //  && fs.existsSync(this.formdata.mdd) && fs.existsSync(this.formdata.css) && fs.existsSync(this.formdata.js)
      //  && this.formdata.mdx.endsWith('.mdx') && this.formdata.js.endsWith('.js') && this.formdata.css.endsWith('.css')
      if (!fs.existsSync(this.formdata.mdd) || !fs.formdata.mdd.endsWith('.mdd')) {
        this.info('mdd未生效', '未选择或者mdd文件非法')
      }
      if (!fs.existsSync(this.formdata.js) || !fs.formdata.js.endsWith('.js')) {
        this.info('js未生效', '未选择或者js文件非法')
      }
      if (!fs.existsSync(this.formdata.css) || !fs.formdata.css.endsWith('.css')) {
        this.info('css未生效', '未选择或者css文件非法')
      }
      if (this.formdata.mdd.endsWith('.mdd')) {
        if (fs.existsSync(this.formdata.mdx)) {
          this.$store.dispatch('updateMDD', this.formdata.mdd)
          this.$store.dispatch('updateMDX', this.formdata.mdx)
          this.$store.dispatch('updateJS', this.formdata.js)
          this.$store.dispatch('updateCSS', this.formdata.css)
          this.success('保存成功', '')
          // TODO restart bgwindow
          ipcRenderer.send('restartBG')
        } else {
          this.failed('保存失败', 'mdx文件不存在')
        }
      } else {
        this.failed('保存失败', '请检查mdx文件路径')
      }
    },
    reset () {
      console.log('reset')
      this.formdata.mdx = this.$store.state.Query.mdx
      this.formdata.mdd = this.$store.state.Query.mdd
      this.formdata.js = this.$store.state.Query.js
      this.formdata.css = this.$store.state.Query.css
      this.info('重置为配置初始值', '')
    },
    success (title, info) {
      this.$Notice.success({
        title: title,
        desc: info
      })
    },
    info (title, info) {
      this.$Notice.info({
        title: title,
        desc: info
      })
    },
    failed (title, info) {
      this.$Notice.error({
        title: title,
        desc: info
      })
    }
  }
}
</script>


