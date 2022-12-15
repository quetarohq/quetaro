drop table if exists jobs cascade;

drop type if exists job_status;

create type job_status as enum (
  'invalid',
  'pending',
  'invoke_failure',
  'invoked',
  'failure',
  'pass'
);

create table jobs (
  id uuid primary key,
  queue_name varchar(64) not null,
  function_name varchar(64) not null,
  payload jsonb not null,
  status job_status not null,
  invoke_after timestamptz not null,
  error_count integer not null,
  last_error text,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

create index if not exists jobs_dequeue_idx on jobs (queue_name, status, invoke_after, updated_at);
