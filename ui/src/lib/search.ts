
import type Comic from '$lib/Comic.svelte';

export type SearchResult = {
	count: number;
	comics: Comic[];
	time: number;
};

function extractDateRange(input: string): string[] {
	var rx = /@date:\s?([0-9]{4}\-[0-9]{2}\-[0-9]{2})\s?[,]?\s?([0-9]{4}\-[0-9]{2}\-[0-9]{2})?/g;
	var regexp = new RegExp(rx);
	return regexp.exec(input);
}

function epoch(dateStr: string): number {
	var ms = new Date(dateStr).getTime();
	return Math.floor(ms / 1000);
}

export async function search(query: string): Promise<SearchResult | null> {
	if (query == "") {
		return null
	}

	query = sanitize(query)

	// handle date query
	if (query.includes("@date")) {
		var match = extractDateRange(query)
		var from = epoch(match[1])

		if (match[2] == undefined) {
			var to = Math.floor(Date.now() / 1000)
		} else {
			var to = epoch(match[2])
		}

		if (!isNaN(from) && !isNaN(to)) {
			query = `@date:[${from} ${to}]`
		} else {
			throw Error("invalid date format")
		}
	}

	let abort = new AbortController();
	try {
		const response = await fetch("/search?q=" + query, {
			signal: abort.signal
		});

		if (!response.ok) {
			const result = await response.json();
			throw Error(`error: ${result.error}`);
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
			return null
		} else if (err.message.includes("Syntax error")) {
			throw Error(`error: Invalid query format`)
		} else {
			console.error(err);
			throw err;
		}
	}
};

