<template>
  <li>
    <span class="text-overflow">
      <button @click="deleteSong" class="delete"><span class="glyphicon glyphicon-trash"></span></button>
      <button :disabled="disableup" @click="up" class="select-song pull-right"><span class="glyphicon glyphicon-arrow-up"></span></button>
      <button :disabled="disabledown" @click="down" class="select-song pull-right" style="margin-left: 15px;"><span class="glyphicon glyphicon-arrow-down"></span></button>
      {{ item.Media.Title }}
    </span>
    <div v-if="item.Downloading" class="progress" style="margin-top: 8px; height: 2vh; font-size: 12px; text-align: center; line-height: 100%;">
      <div class="progress-bar" :class="item.DownloadProgress < 50 ? 'progress-bar-info' : 'progress-bar-success'" style="min-width: 2vw;" :style="{width: item.DownloadProgress + '%'}"></div>
      <p style="text-align: center; width: 100%; position: absolute;">{{ downloadText() }}</p>
    </div>
  </li>
</template>

<script lang="ts">
  import $ from 'jquery'
  import { Component, Prop, Vue } from 'vue-property-decorator'

  @Component
  export default class QueueItem extends Vue {
    @Prop(Boolean) disabledown!: boolean
    @Prop(Boolean) disableup!: boolean
    @Prop(Number) index!: number
    @Prop(Object) item!: IQueueItem

    deleteSong () {
      $.ajax({
        type: 'DELETE',
        url: '/queue/' + this.index
      })
    }

    down () {
      $.ajax({
        data: { move: 'down' },
        type: 'PUT',
        url: '/queue/' + this.index
      })
    }

    downloadText () {
      if (this.item.DownloadProgress < 50) {
        return 'Downloading (' + Math.floor(this.item.DownloadProgress) + '%)'
      } else {
        return 'Transcoding (' + Math.floor(this.item.DownloadProgress) + '%)'
      }
    }

    up () {
      $.ajax({
        data: { move: 'up' },
        type: 'PUT',
        url: '/queue/' + this.index
      })
    }
  }
</script>