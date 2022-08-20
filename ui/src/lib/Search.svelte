<script lang="ts">
	let query: string
</script>

<form class="container search">
	<label for="search-bar" class="sr-only">Search</label>
	<input id="search-bar" type="search"
		autofocus
		placeholder="Search..."
		bind:value={query}
	/>

	{#if query}
		<p>
			Found results in time s.
		</p>
		<div>
			Result 1
		</div>
	{/if}
</form>

<style>
.search {
	width: 75%;
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
