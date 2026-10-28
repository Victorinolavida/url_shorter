--
-- PostgreSQL database dump
--


-- Dumped from database version 17.11
-- Dumped by pg_dump version 17.11

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: url_code; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.url_code (
    id bigint NOT NULL,
    code text NOT NULL,
    url text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: url_code_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.url_code_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: url_code_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.url_code_id_seq OWNED BY public.url_code.id;


--
-- Name: url_code id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.url_code ALTER COLUMN id SET DEFAULT nextval('public.url_code_id_seq'::regclass);


--
-- Name: url_code url_code_code_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.url_code
    ADD CONSTRAINT url_code_code_key UNIQUE (code);


--
-- Name: url_code url_code_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.url_code
    ADD CONSTRAINT url_code_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--


