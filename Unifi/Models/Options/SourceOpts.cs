using System;
using System.ComponentModel.DataAnnotations;

namespace Unifi.Models.Options
{
    /// <summary>
    /// The source options
    /// </summary>
    public record SourceOpts
    {
        public const string Section = "Unifi";

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        [Required(ErrorMessage = Section + ":" + nameof(Host) + " is missing")]
        public string Host { get; init; } = "https://unifi.local:8443";

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        [Required(ErrorMessage = Section + ":" + nameof(Username) + " is missing")]
        public string Username { get; init; } = "unifi";

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        [Required(ErrorMessage = Section + ":" + nameof(Password) + " is missing")]
        public string Password { get; init; } = "unifi";

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        [Required(ErrorMessage = Section + ":" + nameof(Site) + " is missing")]
        public string Site { get; init; } = "default";

        /// <summary>
        /// 
        /// </summary>
        /// <returns></returns>
        [Required(ErrorMessage = Section + ":" + nameof(AwayTimeout) + " is missing")]
        public TimeSpan AwayTimeout { get; init; } = new(0, 5, 1);

        /// <summary>
        /// 
        /// </summary>
        /// <returns></returns>
        [Required(ErrorMessage = Section + ":" + nameof(PollingInterval) + " is missing")]
        public TimeSpan PollingInterval { get; init; } = new(0, 0, 11);
    }
}
