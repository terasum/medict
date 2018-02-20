<template>
<div :class="['logo', globalState.state, globalState.state === 'normal' ? 'center' : 'left' ]">
      <span class="logo-img"></span>
      <h1>MeDict</h1>
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
.logo {
  width: 100vw;
  transition: width 0.1s ease-in-out;
  -webkit-transition: width 0.1s ease-in-out;
}

.center {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: row;
}

.left {
  width: 160px;
  display: flex;
  justify-content: flex-start;
  align-items: center;
  flex-direction: row;
}
/* normal part */

.normal{
  margin-top: 80px;
  margin-bottom: 20px;
  height: 100px;
  user-select: none;
  cursor: default;

  &>h1{
    font-size: 64px;
    line-height: 100px;
    font-weight: lighter;
    font-family: "Monoton";
    color: #fff;
    width: 265px;
    height: 100px;
    cursor: default;

  }

  .logo-img {
    background: url("~@/assets/images/medict.svg");
    background-size: cover;
    width: 54px;
    height: 54px;
  }
}

/* search state style */
.search {
  margin-left: 10px;
  width: 146px;
  height: 36px;
  user-select: none;
  cursor: default;
  & > h1 {
    padding-left: 3px;
    font-size: 28px;
    line-height: 31px;
    font-weight: 300;
    font-family: "Monoton";
    color: #fff;
    width: 146px;
    height: 31px;
    cursor: default;

  }

  .logo-img {
    display: inline-flex;
    background: url("~@/assets/images/medict.svg") 50% 50% no-repeat;
    background-size: cover;
    flex-direction: column;
    width: 25px;
    height: 28px;
  }
}
</style>

