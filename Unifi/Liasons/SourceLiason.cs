using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using TwoMQTT.Interfaces;
using TwoMQTT.Liasons;
using Unifi.DataAccess;
using Unifi.Models.Options;
using Unifi.Models.Shared;
using Unifi.Models.Source;

namespace Unifi.Liasons
{
    /// <summary>
    /// A class representing a managed way to interact with a source.
    /// </summary>
    public class SourceLiason : PollingSourceLiasonBase<Resource, SlugMapping, ISourceDAO, SharedOpts>, ISourceLiason<Resource, object>
    {
        public SourceLiason(ILogger<SourceLiason> logger, ISourceDAO sourceDAO,
            IOptions<SourceOpts> opts, IOptions<SharedOpts> sharedOpts) :
            base(logger, sourceDAO, sharedOpts)
        {
            this.Logger.LogInformation(
                "Host: {host}\n" +
                "Username: {username}\n" +
                "Password: {password}\n" +
                "Site: {site}\n" +
                "AwayTimeout: {awayTimeout}\n" +
                "PollingInterval: {pollingInterval}\n" +
                "AsDeviceTracker: {asDeviceTracker}\n" +
                "Resources: {@resources}\n" +
                "",
                opts.Value.Host,
                opts.Value.Username,
                (!string.IsNullOrEmpty(opts.Value.Password) ? "<REDACTED>" : string.Empty),
                opts.Value.Site,
                opts.Value.AwayTimeout,
                opts.Value.PollingInterval,
                sharedOpts.Value.AsDeviceTracker,
                sharedOpts.Value.Resources
            );
        }

        /// <inheritdoc />
        protected override async Task<Resource?> FetchOneAsync(SlugMapping key, CancellationToken cancellationToken)
        {
            var result = await this.SourceDAO.FetchOneAsync(key, cancellationToken);
            return result switch
            {
                Response => new Resource
                {
                    Mac = result.MACAddress,
                    State = result.State,
                },
                _ => null,
            };
        }
    }
}
