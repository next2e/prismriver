<template>
  <div class="progress" style="margin-top: 8px; height: 2vh; font-size: 12px; text-align: center; line-height: 100%; position: relative;">
    <div v-if="error" class="progress-bar progress-bar-danger" style="width: 100%"></div>
    <div v-else class="progress-bar" :class="progress < 50 ? 'progress-bar-info' : 'progress-bar-success'" style="min-width: 2vw;" :style="{width: progress + '%'}"></div>
    <p style="text-align: center; width: 100%; position: absolute; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ displayText() }}</p>
  </div>
</template>

<script lang="ts">
  import { Component, Prop, Vue } from 'vue-property-decorator'

  @Component
  export default class ProgressBar extends Vue {
    @Prop(String) error!: string
    @Prop(Number) progress!: number

    displayText () {
      if (this.error) {
        return this.error
      } else if (this.progress < 50) {
        return 'Downloading (' + Math.floor(this.progress) + '%)'
      } else {
        return 'Transcoding (' + Math.floor(this.progress) + '%)'
      }
    }
  }
</script>