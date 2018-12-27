<template>
  <li>
    <span>
      <button @click="deleteSong" class="delete"> Delete </button>
      <button v-if="!disableup" @click="up" class="select-song">Up</button>
      <button v-if="!disabledown" @click="down" class="select-song">Down</button>
      {{ title }}
    </span>
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
    @Prop(String) title!: string

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