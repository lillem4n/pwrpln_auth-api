SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
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
-- Name: accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.accounts (
    id uuid NOT NULL,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    name text NOT NULL,
    "apiKey" text,
    password text
);


--
-- Name: accountsFields; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."accountsFields" (
    id uuid NOT NULL,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "accountId" uuid NOT NULL,
    name text NOT NULL,
    value text[] NOT NULL
);


--
-- Name: renewalTokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."renewalTokens" (
    "accountId" uuid NOT NULL,
    exp timestamp without time zone DEFAULT (CURRENT_TIMESTAMP + '24:00:00'::interval) NOT NULL,
    token character(60) NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(255) NOT NULL
);


--
-- Name: accountsFields accountsFields_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."accountsFields"
    ADD CONSTRAINT "accountsFields_pkey" PRIMARY KEY (id);


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: idx_accountname; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_accountname ON public.accounts USING btree (name);


--
-- Name: idx_accountsfields; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_accountsfields ON public."accountsFields" USING btree ("accountId", name);


--
-- Name: idx_renewaltokensaccountid; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_renewaltokensaccountid ON public."renewalTokens" USING btree ("accountId");


--
-- Name: idx_renewaltokensexp; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_renewaltokensexp ON public."renewalTokens" USING btree (exp);


--
-- Name: idx_renewaltokenstoken; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_renewaltokenstoken ON public."renewalTokens" USING btree (token);


--
-- Name: accountsFields accountsFields_accountId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."accountsFields"
    ADD CONSTRAINT "accountsFields_accountId_fkey" FOREIGN KEY ("accountId") REFERENCES public.accounts(id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- Name: renewalTokens renewalTokens_accountId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."renewalTokens"
    ADD CONSTRAINT "renewalTokens_accountId_fkey" FOREIGN KEY ("accountId") REFERENCES public.accounts(id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20201207191913');
