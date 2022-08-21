
import Comic from '$lib/Comic.svelte';

export type SearchResult = {
	count: number;
	comics: Comic[];
	time: number;
};

export async function search(query: string): Promise<SearchResult | null> {
	if (query == "") {
		return null
	}

	let abort = new AbortController();
	try {
		// TODO change to /search?q= when embedded in binary
		const response = await fetch("http://localhost:6380/search?q=" + query, {
			signal: abort.signal
		});

		if (!response.ok) {
			const result = await response.json();
			throw Error(`error ${response.status}: ${result.error}`);
		}
		const result = await response.json();
		return {
			count: result.count,
			comics: result.results,
			time: result.query_time,
		};
	} catch (err: any) {
		if (err.name == 'AbortError') {
			console.log('fetch aborted');
		} else {
			console.error(err);
			throw err;
		}
	}
};

