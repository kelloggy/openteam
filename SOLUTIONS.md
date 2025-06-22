## Solution notes

### Task 01 – Run‑Length Encoder
- Language: Go
- Approach: At first, I planned to use hashmap like with key and value. But my test case were failing and I realized hashmap not restored the order. So I have to switch my plan to use rune which is an int and can also hold special character like emoji. 

I first set the result as empty string, current character as runes[0] and count as zero. I then convert the string into rune and loop that rune. While looping, if the current character is the same as the the loop char, increment the count. If not, concat the result string and reset currentChar as [i] and count to 0. 
- Why: Hashmaps are good for counting but bad for preserving order. Run-length encoding is consecutive runs, not total counts.
- Time spent: ~20 min
- AI tools used: [IF ANY]

### Task 02 – Fix‑the‑Bug
- Language: Go
- Approach: Replaced the non-thread-safe increment (`counter++`) with `atomic.AddInt32(&counter, 1)` to make the ID generation safe under concurrency. I first noticed that the bug was due to multiple goroutines incrementing the same variable without synchronization, which can lead to race conditions and duplicate IDs. I considered using a `sync.Mutex` to lock the increment operation, but realized that for a simple counter, atomic operations are more efficient and idiomatic in Go. With atomic, we can make sure that one instruction will happen one at a time in sequence without any interferance even if we use concurrent goroutine.
- Why: I considered both `sync.Mutex` and `sync/atomic`. Since we are only incrementing a single `int32` variable, `atomic` is the more efficient and minimal solution. It avoids locking overhead and keeps the fix small, fast, and idiomatic.
- Time spent: ~15 min
AI tools used: Consulted ChatGPT to research best practices for safe concurrent increments in Go, and to compare sync.Mutex vs sync/atomic.


### Task 03 – Sync‑Aggregator
- Language: Go
- Approach: I built a concurrent file processor using goroutines and channels. Each worker picks up jobs from a shared channel and processes one file at a time. Using channels helped make sure workers get their own files. 
At first, the tests were failing because I started the timeout timer (context.WithTimeout) too early — before opening the file or handling the sleep directive. That caused even valid files to hit the timeout. So I moved the timeout to start after handling the optional #sleep=N line at the top of the file.
But then I hit another edge case: when sleep == timeout, the test still expected a valid result. I assumed it should timeout, but apparently the timeout only applies to processing, not sleeping. So I added a check — if sleep > timeout, I skip everything and early return. Otherwise, I sleep as requested, then proceed with counting lines and words.
I also wrapped the scanning logic in a goroutine and sent the result through a buffered channel. Using a select lets me handle whichever comes first — the result or a timeout.

- Why: It ensures that each file is processed concurrently, correct timeout contract after sleep, and always returns results in the correct order — even if workers finish at different times.
- Time spent: ~1 hr 30 min
- AI tools used: chatgpt


### Task 04 A – sql-reasoning
Language: SQL
Approach: I start from the campaign table because I want to list every campaign, regardless of whether it has pledges or not. So I used LEFT JOIN with pledge to make sure we include campaigns with zero pledges also.
Then I summed up amount_thb using SUM() and wrapped it with COALESCE so it returns 0 instead of NULL when there are no pledges. I then used formula to get pct_of_target and used ROUND(..., 4) to get 4 decimal places.
Finally, I grouped by campaign_id and target_thb, and ordered by pct_of_target descending, then campaign_id ascending.
- Time spent: 20 min
- AI tools used: chatgpt query syntax

### Task 04 B – sql-reasoning
- Language: sql
- Approach: I decided to only think about one case - thailand. We only want to focus on amount_thb from pledges and country from donar table. The question say we should order ascending by amount_thb and filter out with thailand as country. I still need to put the row number for 90th pecentaile and total row count for the formula ceil(0.9 × N). So i did that with row_number() and put the total number of rows as total and dump all the data as ranked. I execute the query in my local db and it execute correctly.
But when i paste my query in this file, the test cas failed. it say no such function: CEIL. I use chatgpt to search for what other func i can use in sqlite. it say to use cast. So i decided to use cast. I use the same method for the global case also. 
- Why:  
- Time spent: 30 min
- AI tools used: chatgpt query syntax

### Task 04 C – sql-reasoning
- Language: sql
- Approach: I decided to create indexes on columns that are heavily used in filtering, joining, and sorting operations to improve query performance.
I added an index on donor.country because it is frequently used as a filter condition which helps narrow down rows quickly.
I also indexed pledge.amount_thb since it is used in ORDER BY and percentile calculations, making sorting and row lookup more efficient.
I did not add indexes on id fields because primary keys in SQLite already have indexes by default.
- Why:  
- Time spent: 5 min
- AI tools used: