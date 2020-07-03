<template>
  <div class="nav">
    <router-link to="/">Jobs</router-link> |
    <router-link to="/ansible">Ansible</router-link> |
    <router-link to="/blueprint">Blueprint</router-link> |
    <router-link to="/load">Load</router-link> |
    <router-link to="/about">About</router-link>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'

@Component
export default class Nav extends Vue {
}
</script>

<style scoped lang="stylus">
.nav
  padding 30px
  border-bottom 1px solid #aaa
  a
    font-weight bold
    color #2c3e50
    &.router-link-exact-active
      color #42b983
</style>
