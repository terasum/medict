<template>
    <div class="layout" :style="{ height: '100%', overflow:'hidden'}">
        <Layout :style="{ height: '100%', overflow:'hidden'}">
            <div class="head-placeholder"></div>
            <Header>
               <div class="header-content">
                   <div class="header-logo header-item"></div>
                   <Input class="header-item" :value="word" @input="searchWord" icon="search" placeholder="Enter something..." style="width: 200px"/>
                   <ButtonGroup>
                     <Button icon="arrow-left-b"></Button>
                     <Button icon="arrow-right-b"></Button>
                   </ButtonGroup>
                   <Select v-model="model1" style="width:200px">
                     <Option v-for="item in cityList" :value="item.value" :key="item.value">{{ item.label }}</Option>
                   </Select>
                   <Dropdown trigger="click">
                        <Button icon="settings"></Button>
                        <DropdownMenu slot="list">
                            <DropdownItem>界面设置</DropdownItem>
                            <DropdownItem>词典设置</DropdownItem>
                            <DropdownItem divided>关于</DropdownItem>
                        </DropdownMenu>
                    </Dropdown>
               </div> 
            </Header>
            <Layout :style="{ overflow: 'hidden', background: '#fff', height: '100%'}">
                <Sider hide-trigger :style="{ overflow: 'hidden', background: '#fff', height: '100%'}">
                    
                </Sider>
                <Layout>
                    <Content class="dict-content">
                        <!-- {{cwdpath}}  -->
                        {{content}}
                    </Content>
                </Layout>
            </Layout>
        </Layout>
    </div>
</template>


<style scoped>

.head-placeholder{
    -webkit-app-region: drag;
    width: 100%;
    height: 24px;
    background: #0073F1;
    /* background: #f5f7f9; */

}
.layout{
    background: #f5f7f9;
    position: relative;
    overflow: hidden;
}
.header-item{
    margin-left: 5px;
}
.header-logo{
    width: 120px;
    height: 32px;
    background: #f9f7f5;
    border-radius: 3px;
    position: relative;
    vertical-align: middle;
    line-height: normal;
    display: inline-block;
    line-height: 36px;
}
.layout-nav{
    width: 420px;
    margin: 0 auto;
    margin-right: 20px;
}

.layout-content{
padding: '0 2px' ;
 overflowY: 'auto'; 
 width:'100%';
 overflowX:'hidden'
}

.dict-content{
 padding: '2px';
 minHeight: '280px'; 
 height:'100%';
 width:'100%';
 maxWidth:'541px';
 ackground: '#fff';
 overflowWrap: 'break-word';
 overflowY:'auto';
 overflowX:'hidden';
}

#content::-webkit-scrollbar-track
{
	-webkit-box-shadow: inset 0 0 6px rgba(0,0,0,0.3);
	border-radius: 10px;
	background-color: #F5F5F5;
}

#content::-webkit-scrollbar
{
	width: 12px;
	background-color: #F5F5F5;
}

#content::-webkit-scrollbar-thumb
{
	border-radius: 10px;
	-webkit-box-shadow: inset 0 0 6px rgba(0,0,0,.3);
	background-color: #555;
}
</style>


<script>
    export default {
      data () {
        return {
          value4: '',
          cityList: [
            {
              value: 'New York',
              label: 'New York'
            },
            {
              value: 'London',
              label: 'London'
            }
          ],
          model1: ''
        }
      },
      computed: {
        // 当前搜索词
        word () {
          return this.$store.state.medict.searchWord
        },
        cwdpath () {
          return this.$store.state.medict
        },
        // 当前搜索解释
        content () {
          return this.$store.state.medict.content
        }
      },
      methods: {
        searchWord (word) {
          console.log(word)
          console.log(this)
          this.$store.dispatch('init')
          this.$store.dispatch('search', word)
        }
      },
      watch: {
        // 观察变化
        // myvalue: (newVal, oldVal)=>  { }
      }
    }
</script>