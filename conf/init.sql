--  This file is part of the eliona project.
--  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
--  ______ _ _
-- |  ____| (_)
-- | |__  | |_  ___  _ __   __ _
-- |  __| | | |/ _ \| '_ \ / _` |
-- | |____| | | (_) | | | | (_| |
-- |______|_|_|\___/|_| |_|\__,_|
--
--  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
--  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
--  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
--  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
--  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

create schema if not exists thingdust;

--
-- Todo: create tables and database objects necessary for this app like tables persisting configuration
--
-- .
-- This table should be made editable by eliona frontend.
create table if not exists thingdust.config
(
    config_id           bigserial primary key,
    api_endpoint        text not null,
    api_key             text not null,
    enable              boolean   default false,
    refresh_interval    integer default 60,
    request_timeout     integer default 120,
    active              boolean default false,
    project_ids            text[]
);

create table if not exists thingdust.spaces
(
    config_id       bigint not null,
    project_id      text not null,
    asset_id        integer not null,
    space_name      text not null,
    primary key(config_id, project_id, asset_id, space_name)
);


-- Makes the new objects available for all other init steps
commit;


