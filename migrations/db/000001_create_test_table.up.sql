CREATE TABLE IF NOT EXISTS public.test_table (
	id_test uuid NOT NULL,
	desc_test varchar NULL,
	CONSTRAINT test_table_pk PRIMARY KEY (id_test)
);
