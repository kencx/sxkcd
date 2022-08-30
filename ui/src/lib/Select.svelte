<script lang="ts">
	export let selected: string = 'relevancy';
</script>

<div class="select">
	<select bind:value={selected}>
		<option value="relevancy">Relevancy</option>
		<option value="oldest">Least Recent</option>
		<option value="newest">Most Recent</option>
	</select>
</div>

<style>
.select {
	width: 165px;
}

.select select{
	font-size: 0.7rem;
}

@media (min-width:320px) and (max-width:640px) {
	.select {
		width: 60%;
	}
}
</style>
