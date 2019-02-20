<template>
  <div class="component" style="padding-top: 6px;">
    <p class="linkP"> Choose from past songs &nbsp; <a href="#search-container" id="toggle-search" data-toggle="collapse" aria-expanded="false" aria-controls="search-container" @click="toggle"><span id="search-arrow" class="glyphicon glyphicon-chevron-down"></span></a></p>
    <div id="search-container" class="collapse"> <!-- had style="display: none;" -->
      <div class="linkContainer">
        <form class="linkForm">
          <input id="query" class="linkInput" type="text" autocomplete="off" autofocus placeholder="Search" v-model="query" @keydown.enter.prevent="submit">
          <button id="submit-search" class="linkButton hvr-shutter-out-horizontal" type="button" @click="submit">Search</button>
        </form>
        <button id="shuffle" class="hvr-shutter-out-horizontal" type="button" @click="shuffle">Shuffle</button>
      </div>
      <ul id="search-results">
        <search-item v-for="result in results" :id="result.ID" :title="result.Title" :type="result.Type"></search-item>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
  import $ from 'jquery'
  import { Component , Vue } from 'vue-property-decorator'

  @Component
  export default class SearchForm extends Vue {
    query = ''
    results: IMedia[] = []

    getRandomSearch () {
      const params = new URLSearchParams()
      params.set('limit', '20')
      fetch('media/random?' + params.toString()).then((response) => {
        return response.json()
      }).then((json) => {
        this.results = json as IMedia[]
      })
    }

    mounted () {
      this.getRandomSearch()
    }

    shuffle (event: Event) {
      $((event.target as Object)).blur()
      this.getRandomSearch()
    }

    submit (): void {
      if (this.query.length) {
        const params = new URLSearchParams()
        params.set('query', this.query)
        fetch('media/search?' + params.toString()).then((response) => {
          return response.json()
        }).then((json) => {
          this.results = json as IMedia[]
        })
      } else {
        this.getRandomSearch()
      }
    }

    toggle (event: Event): void {
      $((event.target as Object)).blur()
      $((event.target as Object)).toggleClass('glyphicon-chevron-down')
      $((event.target as Object)).toggleClass('glyphicon-chevron-up')
    }
  }
</script>