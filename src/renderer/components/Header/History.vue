<template>
  <div :class="['history', globalState.state]">
    <div class="navi-btn prev"></div>
    <div class="navi-btn next"></div>
  </div>
</template>

<script>
export default {
  data () {
    return {
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
  flex-direction: row;
  justify-content: flex-start;
  align-items: center;

  .navi-btn{
    display: flex;
    flex-direction:column;
    justify-content: center;
    width: 30px;
    height: 30px;
    border-radius: 3px;
    margin-left: 10px;
    background: #fff;
    cursor: pointer;
    &:hover{
        background-color: #efefef;
        outline: none;
        cursor: pointer;
    }
    &:active{
        background-color: #ccc;
        box-shadow: 2px 2px #666;
        transform: translate(1px 1px);
        outline: none;
    }
  }
  .prev{
    background: url('~@/assets/images/left.svg') no-repeat center;
    background-color: #fff;
    background-size: 24px 24px;
  }
  .next{
    background: url('~@/assets/images/right.svg') no-repeat center;
    background-color: #fff;
    background-size: 24px 24px;
  }
}

</style>



