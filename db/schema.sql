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
-- Name: links; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.links (
    id bigint NOT NULL,
    code text NOT NULL,
    raw_url text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    url_hash bytea NOT NULL,
    expires_at timestamp with time zone,
    status text DEFAULT 'active'::text NOT NULL,
    CONSTRAINT chck_code_length CHECK (((char_length(code) >= 4) AND (char_length(code) <= 32))),
    CONSTRAINT chck_raw_url_length CHECK ((length(raw_url) <= 2048)),
    CONSTRAINT chck_status CHECK ((status = ANY (ARRAY['active'::text, 'expired'::text, 'disabled'::text, 'blocked'::text])))
);


--
-- Name: links_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.links ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.links_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: links unique_link_code; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.links
    ADD CONSTRAINT unique_link_code UNIQUE (code);


--
-- Name: links url_code_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.links
    ADD CONSTRAINT url_code_pkey PRIMARY KEY (id);


--
-- Name: idx_link_expired_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_link_expired_at ON public.links USING btree (expires_at) WHERE (expires_at IS NOT NULL);


--
-- Name: idx_url_hash_expired_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_url_hash_expired_at ON public.links USING btree (url_hash);


--
-- PostgreSQL database dump complete
--


