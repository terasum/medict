<template>
  <div :class="['select-dict', global.state ]">
  </div>
</template>

<script>
export default {
  data () {
    return {
      selected: null
    }
  },
  props: {
    global: {
      type: Object,
      required: true
    }
  }
}
</script>

<style lang="scss" scoped>
.normal{
  display: none;
}

.search{
  display: flex;
  background-color: #fff;
  border: none;
  border-radius: 4px;  
  width:161px;
  height: 30px;
  margin-left: 10px;
}
</style>



