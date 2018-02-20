<template>
  <div :class="['select-dict', globalState.state ]">
    <span>{{globalState.dict.name}}</span> 
  </div>
</template>

<script>
export default {
  data () {
    return {
      selected: null,
      globalState: this.$store.state.MedictGlobal
    }
  }
}
</script>

<style lang="scss" scoped>
.normal {
  display: none;
}

.search {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #fff;
  border: none;
  border-radius: 4px;
  width: 161px;
  height: 30px;
  margin-left: 10px;
  cursor: pointer;
  & > span {
    color: #1377ed;
    line-height: 30px;
    font-size: 14px;
    user-select: none;
    &:before {
      padding-right: 4px;
      font-family: fontAwesome;
      content: "\f02d";
    }
  }
  &:hover {
    background: #efefef;
    outline: none;
    cursor: pointer;
  }
  &:active {
    background: #ccc;
    box-shadow: 2px 2px #666;
    transform: translate(1px 1px);
    outline: none;
  }
}
</style>



