using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.CompilerServices;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using TwoMQTT.Core.Interfaces;
using Unifi.DataAccess;
using Unifi.Models.Shared;
using Unifi.Models.Source;

namespace Unifi.Liasons
{
    /// <summary>
    /// A class representing a managed way to interact with a source.
    /// </summary>
    public class SourceLiason : ISourceLiason<Resource, Command>
    {
        public SourceLiason(ILogger<SourceLiason> logger, ISourceDAO sourceDAO,
            IOptions<Models.Options.SourceOpts> opts, IOptions<Models.Options.SharedOpts> sharedOpts)
        {
            this.Logger = logger;
            this.SourceDAO = sourceDAO;
            this.Questions = sharedOpts.Value.Resources;

            this.Logger.LogInformation(
                "Host: {host}\n" +
                "Username: {username}\n" +
                "Password: {password}\n" +
                "Site: {site}\n" +
                "AwayTimeout: {awayTimeout}\n" +
                "PollingInterval: {pollingInterval}\n" +
                "Resources: {@resources}\n" +
                "",
                opts.Value.Host,
                opts.Value.Username,
                (!string.IsNullOrEmpty(opts.Value.Password) ? "<REDACTED>" : string.Empty),
                opts.Value.Site,
                opts.Value.AwayTimeout,
                opts.Value.PollingInterval,
                sharedOpts.Value.Resources
            );
        }

        /// <inheritdoc />
        public async IAsyncEnumerable<Resource?> FetchAllAsync([EnumeratorCancellation] CancellationToken cancellationToken = default)
        {
            foreach (var key in this.Questions)
            {
                this.Logger.LogDebug("Looking up {key}", key);
                var result = await this.SourceDAO.FetchOneAsync(key, cancellationToken);
                var resp = result != null ? this.MapData(result) : null;
                yield return resp;
            }
        }

        /// <summary>
        /// The logger used internally.
        /// </summary>
        private readonly ILogger<SourceLiason> Logger;

        /// <summary>
        /// The dao used to interact with the source.
        /// </summary>
        private readonly ISourceDAO SourceDAO;

        /// <summary>
        /// The questions to ask the source (typically some kind of key/slug pairing).
        /// </summary>
        private readonly List<SlugMapping> Questions;

        /// <summary>
        /// Map the source response to a shared response representation.
        /// </summary>
        /// <param name="src"></param>
        /// <returns></returns>
        private Resource MapData(FetchResponse src) =>
            new Resource
            {
                Mac = src.MACAddress,
                State = src.State,
            };
    }
}
