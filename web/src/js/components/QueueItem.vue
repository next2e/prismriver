<template>
  <li>
    <span>
      <button @click="deleteSong" class="delete"> Delete </button>
      {{ title }}
    </span>
  </li>
</template>

<script lang="ts">
  import $ from 'jquery'
  import { Component, Prop, Vue } from 'vue-property-decorator'

  @Component
  export default class QueueItem extends Vue {
    @Prop(Number) index!: number
    @Prop(String) title!: string

    deleteSong () {
      $.ajax({
        type: 'DELETE',
        url: '/queue/' + this.index
      })
    }
  }
</script>