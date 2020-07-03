<template>
  <div class="home">
    <div v-for="(job, key) in $store.state.jobs" :key="key">
      {{ job }}
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'

@Component({
  components: {},
  mounted() {
    this.$store.dispatch('getJobs')
  },
})
export default class Home extends Vue {}
</script>
