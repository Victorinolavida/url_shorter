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
-- Name: click_events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.click_events (
    id bigint NOT NULL,
    link_id bigint NOT NULL,
    clicked_at timestamp with time zone DEFAULT now() NOT NULL,
    user_agent text DEFAULT ''::text NOT NULL,
    browser text DEFAULT ''::text NOT NULL,
    os text DEFAULT ''::text NOT NULL,
    device_type text DEFAULT ''::text NOT NULL,
    country text DEFAULT ''::text NOT NULL,
    ip_hash bytea,
    is_bot boolean DEFAULT false NOT NULL,
    status_code smallint NOT NULL,
    CONSTRAINT chck_click_events_status_code CHECK ((status_code = ANY (ARRAY[301, 303, 302, 410]))),
    CONSTRAINT chck_click_events_user_agent_length CHECK ((char_length(user_agent) <= 512))
);


--
-- Name: click_events_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.click_events ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.click_events_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


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
-- Name: click_events pk_click_events; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.click_events
    ADD CONSTRAINT pk_click_events PRIMARY KEY (id);


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
-- Name: unique_active_url_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_active_url_hash ON public.links USING btree (url_hash) WHERE (status = 'active'::text);


--
-- PostgreSQL database dump complete
--


