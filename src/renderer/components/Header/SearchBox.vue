<template>
    <div :class="[ 'search-box', globalState.state , globalState.state === 'normal' ? 'center': 'left' ]">
        <input 
        type="text" 
        autofocus="autofocus"
        @keyup="predict"
        @keyup.enter="search"
        placeholder="Search some words" >
        <button> {{ searchButtonContent }}</button>
    </div>
</template>

<script>
export default {
  computed: {
    searchButtonContent () {
      return this.globalState.state === 'normal' ? 'Search' : ''
    }
  },
  data () {
    return {
      globalState: this.$store.state.MedictGlobal
    }
  },
  methods: {
    predict (e) {
      let queryWord = e.target.value.trim()
      if (this.globalState.state === 'normal') {
        if (queryWord !== '') {
          this.$store.commit('TRANS_TO_SEARCH', e.target.value)
        }
      }
      let mdict = this.$store.dispatch('predict')
      let _this = this
      mdict.then((_mdict) => {
        const prefix = _mdict.prefix(queryWord)
        _this.$store.dispatch('searchWold', prefix)
      })
    },
    search (e) {
      if (this.globalState.state === 'search') {
        if (e.target.value.trim() === '') {
          this.$store.commit('TRANS_TO_NORMAL')
        }
      }
    }
  }
}
</script>

<style lang="scss" scoped>

.search-box{
    user-select: none;
    transition: width 0.4s ease-in-out;
    -webkit-transition: width 0.4s ease-in-out;
}

.center{
    display: flex;
    flex-direction: row;
    justify-content: center;
    transition: width 0.4s ease-in-out;
    -webkit-transition: width 0.4s ease-in-out;
}

.left{
    display: flex;
    flex-direction: row;
    justify-content: flex-start;
    align-items: center;
    transition: width 0.4s ease-in-out;
    -webkit-transition: width 0.4s ease-in-out;
}

/* normal part */
.normal{
    &>input[type=text]{
        border: 1px #fff solid;
        border-radius: 12px;
        height: 36px;
        width: 357px;
        background-color: #5CA2EF;
        padding-left: 36px;
        background: url("~@/assets/images/search.svg") no-repeat;
        background-color: #5CA2EF;
        background-size: 22px 22px;
        background-position: 6px 6px;

        color:#fff;
        font-size: 16px;
        font-weight: lighter;
        line-height: 36px;

    }
    &>input[type=text]:focus{
        border: 1px #ccc solid;
        outline: none;
        user-select: none;
    }

    &> input::-webkit-input-placeholder{
      color: #f0f0f0;
      font-size: 16px;
      font-weight: lighter;
      line-height: 36px;
      white-space: pre;
    }
    &>button{
        display: none;
    }
}


.search{
    &>input[type=text]{
        border: 1px #fff solid;
        border-radius: 12px;
        height: 30px;
        width: 247px;
        background-color: #5CA2EF;
        padding-left: 10px;
        margin-left: 10px;
        // background: url("~@/assets/images/search.svg") no-repeat;
        background-color: #5CA2EF;
        background-size: 16px 16px;
        background-position: 6px 6px;
        transition: width 0.4s ease-in-out;
        -webkit-transition: width 0.4s ease-in-out;
        /* fonts */
        font-weight: lighter;
        color: #fff;
        font-size: 14px;
        line-height: 30px;

    }

    &>input[type=text]:focus{
        border: 1px #ccc solid;
        outline: none;
        user-select: none;
    }

    &> input::-webkit-input-placeholder{
      color: #f0f0f0;
      font-size: 14px;
      font-weight: lighter;
      line-height: 30px;
      white-space: pre;
    }

    &>button{
        height: 30px;
        width: 30px;
        border-radius: 5px;
        margin-left: 5px;
        border: none;
        background: #fff;
        color:#0073F1;
        font-size: 20px;
        line-height: 30px;
        font-weight: lighter;
        transition: width 0.4s ease-in-out;
        -webkit-transition: width 0.4s ease-in-out;
        user-select: none;
        outline: none;
    }
    &>button:before{
        font-family: FontAwesome;
        content: '\f002';
    }
    &>button:hover{
        background: #EFEFEF;
        outline: none;
        cursor: pointer;
    }
    &>button:active{
        background: #ccc;
        box-shadow: 2px 2px #666;
        transform: translate(1px 1px);
        outline: none;
    }
    
}
</style>
