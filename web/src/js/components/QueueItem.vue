<template>
  <li>
    <span>
      <button @click="deleteSong" class="delete"> Delete </button>
      <button v-if="!disableup" @click="up" class="select-song">Up</button>
      <button v-if="!disabledown" @click="down" class="select-song">Down</button>
      {{ item.Media.Title }}
    </span>
    <div v-if="item.Downloading" class="progress" style="margin-top: 8px; height: 2vh; font-size: 12px; text-align: center; line-height: 100%;">
      <div class="progress-bar" :class="item.DownloadProgress < 50 ? 'progress-bar-info' : 'progress-bar-success'" style="min-width: 2vw;" :style="'width: ' + item.DownloadProgress + '%'"></div>
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

    up () {
      $.ajax({
        data: { move: 'up' },
        type: 'PUT',
        url: '/queue/' + this.index
      })
    }
  }
</script>