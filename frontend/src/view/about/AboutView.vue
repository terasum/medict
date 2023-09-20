<!--

 Copyright (C) 2023 Quan Chen <chenquan_act@163.com>

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
-->

<template>
  <div class="container about-view">
    <div class="about-title">
      <h1>About</h1>
    </div>
    <div class="about-body">
      <section class="about-section">
        <div class="section-title">
          <h2> <NIcon><Code/></NIcon> 软件信息</h2>
        </div>
        <div class="section-body">
          <div class="input-group">
            <div class="input-info">
              <label>软件简介</label>
              <p>Medict, 现代跨平台词典App</p>
            </div>
            <div class="input-container">
            <button
                class="btn"
                @click="openExternal('https://github.com/terasum/medict')"
                >访问 github</button
              >
            </div>
          </div>
          <div class="input-group">
            <div class="input-info">
              <label>软件版本</label>
              <p>v3.0.1-alpha</p>
            </div>
            <div class="input-container">
            </div>
          </div>
        </div>
             </section>

      <section class="about-section">
        <div class="section-title">
          <h2><NIcon><PeopleCarry/></NIcon> 致谢</h2>
        </div>
        <div class="section-body">
          <div class="input-group">
            <div class="input-info">
              <label>开发人员</label>
              <p>
                <a href="javascript:void(0);" @click="openExternal('https://github.com/terasum')">Chen, Quan</a>
              </p>
            </div>
            <div class="input-container">
            </div>
          </div>

          <div class="input-group">
            <div class="input-info">
              <label>设计人员</label>
              <p>Zhang, Mingjiao</p>
            </div>
            <div class="input-container">
            </div>
          </div>
          <div class="input-group">
            <div class="input-info">
              <label>特别鸣谢</label>
              <p>Song, Xing</p>
            </div>
            <div class="input-container">
            </div>
          </div>
        </div>
      </section>

    </div>
  </div>
</template>

<script lang="ts" setup>
import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime';
import {Code, PeopleCarry} from "@vicons/fa";
import { NIcon } from 'naive-ui';

function openExternal(url: string) {
      BrowserOpenURL(url);
}
</script>



<style lang="scss" scoped>
.container{
  height: 100%;
  overflow-y: auto;
}

.about-title {
  padding: 10px 0 10px 12px;
}


.about-section {
  display: flex;
  flex-direction: column;
  margin: 0 10px;

  .section-title {
    width: 100%;
    padding: 2px 4px;
    margin: 0;
    font-size: 16px;
    h2{
      font-size: 21px;
      margin: 14px 0;
      padding: 6px 0;
      border-bottom: 1px solid #c1c1c3;
    }
  }

  .section-body {
    padding-left: 12px;
  }

  .input-group {
    width: 100%;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    margin: 10px 0;
    .input-info {
      label {
        font-style: bold;
      }

      p {
        color: #333;
        font-style: italic;
        margin: 0;
      }
      
      a{
        color: #333;
      }

    }

    .input-container {
      display: flex;
      flex-direction: column;
      justify-content: center;
      max-width: 15rem;
      padding-right: 20px;
      input {
        height: 26px;
      }
    }
  }
}
.about-body {
  min-height: 260px;
}

.btn {
  background-color: rgb(252, 252, 252);
  background-image: linear-gradient(to bottom, #fcfcfc 0, #f1f1f1 100%);
  border-radius: 4px;
  border: 1px solid transparent;
  box-shadow: 0 1px 1px rgb(0 0 0 / 6%);
  border-color: #c2c0c2 #c2c0c2 #a19fa1;
  text-align: center;
  color: #666;
  font-size: 12px;
  display: inline-block;
  padding: 3px 8px;
  margin-bottom: 0;
  white-space: nowrap;
  height: 26px;
  min-width: 80px;
  &:active {
    background-color: #ddd;
    background-image: none;
  }
}



.contributors{
  margin-top: 20px;
}
.contributor {
  margin-top: 10px;
  display: flex;
  justify-content: center;
.contributor-user {
  margin: 5px 8px;
}
}
.contributor-infos{
  padding: 0;
  margin:0;
  border:none;
  background: #f1f1f1;
  user-select: none;

.nickname{
  display: block;
  text-align: center;
  color: #272727;
  font-weight: 400;
  margin:0;
  padding:0;
}

.info{
  display: block;
  text-align: center;
  color: #272727;
  font-weight: lighter;
  margin:0;
  padding:0;
}

}
</style>
