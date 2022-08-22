<script lang="ts">
	import Comic from '$lib/Comic.svelte';
	import type { SearchResult } from './search'
	import { search} from './search'

	let query: string = '';
	let timer;

	const debounce = e => {
		clearTimeout(timer);
		timer = setTimeout(() => {
			query = e.target.value;
		}, 300);
	}

	$: promise = search(query);

	function focus() {
		document.getElementById("search-bar").focus();
	}

	function handleKeyDown(event) {
		switch (event.key) {
			case "/":
				event.preventDefault()
				focus();
				break;
		}
	}
</script>

<svelte:window on:keydown={handleKeyDown}/>

<div class="container search">
	<label for="search-bar" class="sr-only">Search</label>
	<input id="search-bar" type="search"
		autocomplete="off"
		placeholder="Search..."
		on:input|preventDefault={debounce}
	/>

	{#await promise then data}
		{#if data != null}
			<div class="timer">
				<p>found <span class="contrast">{data.count}</span> results in {(data.time*1000).toFixed(3)}ms</p>
			</div>
			<div>
				{#if data.comics}
					{#each data.comics as comic}
					<Comic result={comic}/>
					{/each}
				{/if}
			</div>
		{/if}
	{:catch error}
		<small class="error">An error occured: {error.message}</small>
	{/await}
</div>

<style>
.search {
	width: 75%;
}

.timer {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	margin-bottom: 1.5rem;
}

.timer p {
	font-size: 0.8rem;
}

.contrast {
	color: orange;
}

.error {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	color: red;
	margin-top: 0.5rem;
}

.sr-only {
	position: absolute;
	width: 1px;
	height: 1px;
	padding: 0;
	margin: -1px;
	overflow: hidden;
	clip: rect(0,0,0);
	white-space: nowrap;
	border-width: 0;
}
</style>
